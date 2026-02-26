# 1. CARGA DE VARIABLES (.env)
ifneq ("$(wildcard .env)","")
    include .env
    export
endif

# 2. VARIABLES (Valores del .env o por defecto)
BINARY ?= notasApp
DOCKER_USER   ?= dinotaurent
VERSION       ?= latest
REGISTRY       = docker.io
K8S_YAML      ?= archivoDeploy.yaml

# Nombre de imagen unificado para evitar errores
IMAGE_FULL     = $(REGISTRY)/$(DOCKER_USER)/$(SERVICE_NAME):$(VERSION)

# --- COMANDOS LOCALES (Podman Compose) ---

up:
	@echo "Levantando local con Podman Compose..."
	podman compose up -d

up_build:build_notas-app-go
	@echo "Reconstruyendo y levantando local..."
	podman compose down
	#podman compose up --build -d                         # Para crear sin el mongoweb-notas
	podman compose --profile debug up --build -d          # Para crear el mongoweb-notas tambien

down:
	@echo "Deteniendo contenedores locales..."
	podman compose down

build_notas-app-go:
	@echo "Compilando binario de Go para Linux..."
	env GOOS=linux CGO_ENABLED=0 go build -o ${BINARY} ./cmd/api

# --- COMANDOS DE KUBERNETES ---

k8_apply:
	@echo "Aplicando archivos YAML en el cluster..."
	kubectl apply -f $(K8S_YAML)

k8_delete:
	@echo "Eliminando recursos del cluster..."
	kubectl delete -f $(K8S_YAML)

k8_mongo_deploy:
	@echo "Aplicando archivos YAML  de mongo en el cluster..."
	kubectl apply -f mongoDeploy.yaml

k8_mongo_delete:
	@echo "Eliminando recursos del cluster..."
	kubectl delete -f mongoDeploy.yaml

k8_deploy:
	# Mongo db y express
	@echo "Aplicando archivos YAML  de mongo en el cluster..."
	kubectl apply -f mongoDeploy.yaml

	# Notas App
	@echo "Preparando imagen: $(IMAGE_FULL)..."
	podman build -t $(IMAGE_FULL) .
	@echo "Subiendo a Docker Hub..."
	podman push $(IMAGE_FULL)
	@echo "Actualizando Kubernetes..."
	# Sincroniza el YAML con la versión actual antes de aplicar
	sed -i 's|image: .*|image: $(IMAGE_FULL)|' $(K8S_YAML)
	kubectl apply -f $(K8S_YAML)
	@echo "Forzando reinicio del deployment..."
	kubectl rollout restart deployment/$(SERVICE_NAME)

k8_logs:
	@echo "Mostrando logs del broker..."
	kubectl logs -f deployment/$(SERVICE_NAME)

