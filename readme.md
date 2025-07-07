# File-Binder Makefile
# Ce Makefile automatise la construction et les tests du projet file-binder

# Variables
BINDER_DIR = binder
EXAMPLES_DIR = examples
BUILD_DIR = build
BINDER_BIN = $(BUILD_DIR)/binder
GO_FILES = $(shell find . -name "*.go" -not -path "./$(BUILD_DIR)/*")

# Couleurs pour l'affichage
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[1;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

.PHONY: all build clean test examples help install deps

# Cible par défaut
all: build

# Afficher l'aide
help:
	@echo "$(BLUE)File-Binder - Combinateur d'exécutables$(NC)"
	@echo ""
	@echo "$(YELLOW)Cibles disponibles :$(NC)"
	@echo "  $(GREEN)build$(NC)     - Compiler le binder"
	@echo "  $(GREEN)clean$(NC)     - Nettoyer les fichiers de build"
	@echo "  $(GREEN)test$(NC)      - Exécuter les tests"
	@echo "  $(GREEN)examples$(NC)  - Créer des programmes d'exemple"
	@echo "  $(GREEN)demo$(NC)      - Démonstration complète"
	@echo "  $(GREEN)install$(NC)   - Installer le binder dans \$$GOPATH/bin"
	@echo "  $(GREEN)deps$(NC)      - Vérifier les dépendances"
	@echo "  $(GREEN)help$(NC)      - Afficher cette aide"

# Créer le répertoire de build
$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

# Vérifier les dépendances
deps:
	@echo "$(BLUE)Vérification des dépendances...$(NC)"
	@go version || (echo "$(RED)Go n'est pas installé$(NC)" && exit 1)
	@echo "$(GREEN)✓ Go est installé$(NC)"

# Compiler le binder
build: deps $(BUILD_DIR)
	@echo "$(BLUE)Compilation du binder...$(NC)"
	@cd $(BINDER_DIR) && go build -o ../$(BINDER_BIN) .
	@echo "$(GREEN)✓ Binder compilé : $(BINDER_BIN)$(NC)"

# Installer le binder dans $GOPATH/bin
install: build
	@echo "$(BLUE)Installation du binder...$(NC)"
	@go install ./$(BINDER_DIR)
	@echo "$(GREEN)✓ Binder installé dans \$$GOPATH/bin$(NC)"

# Nettoyer les fichiers de build
clean:
	@echo "$(YELLOW)Nettoyage...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(EXAMPLES_DIR)
	@rm -f combined-*
	@echo "$(GREEN)✓ Nettoyage terminé$(NC)"

# Créer des programmes d'exemple
examples: $(BUILD_DIR)
	@echo "$(BLUE)Création des programmes d'exemple...$(NC)"
	@mkdir -p $(EXAMPLES_DIR)
	
	# Programme A - Hello World
	@echo 'package main' > $(EXAMPLES_DIR)/hello.go
	@echo 'import "fmt"' >> $(EXAMPLES_DIR)/hello.go
	@echo 'func main() {' >> $(EXAMPLES_DIR)/hello.go
	@echo '    fmt.Println("Hello from Program A!")' >> $(EXAMPLES_DIR)/hello.go
	@echo '}' >> $(EXAMPLES_DIR)/hello.go
	@cd $(EXAMPLES_DIR) && go mod init hello && go build -o hello hello.go
	
	# Programme B - Date/Time
	@echo 'package main' > $(EXAMPLES_DIR)/datetime.go
	@echo 'import (' >> $(EXAMPLES_DIR)/datetime.go
	@echo '    "fmt"' >> $(EXAMPLES_DIR)/datetime.go
	@echo '    "time"' >> $(EXAMPLES_DIR)/datetime.go
	@echo ')' >> $(EXAMPLES_DIR)/datetime.go
	@echo 'func main() {' >> $(EXAMPLES_DIR)/datetime.go
	@echo '    fmt.Printf("Current time: %s\n", time.Now().Format("2006-01-02 15:04:05"))' >> $(EXAMPLES_DIR)/datetime.go
	@echo '}' >> $(EXAMPLES_DIR)/datetime.go
	@cd $(EXAMPLES_DIR) && go build -o datetime datetime.go
	
	@echo "$(GREEN)✓ Programmes d'exemple créés dans $(EXAMPLES_DIR)/$(NC)"

# Exécuter les tests
test: build examples
	@echo "$(BLUE)Exécution des tests...$(NC)"
	@echo "$(YELLOW)Test 1: Combinaison des programmes d'exemple$(NC)"
	@./$(BINDER_BIN) $(BUILD_DIR)/combined-test $(EXAMPLES_DIR)/hello $(EXAMPLES_DIR)/datetime
	@echo "$(YELLOW)Test 2: Exécution du programme combiné$(NC)"
	@./$(BUILD_DIR)/combined-test
	@echo "$(GREEN)✓ Tests terminés avec succès$(NC)"

# Démonstration complète
demo: clean build examples test
	@echo "$(GREEN)========================================$(NC)"
	@echo "$(GREEN)  Démonstration File-Binder terminée!$(NC)"
	@echo "$(GREEN)========================================$(NC)"
	@echo ""
	@echo "$(BLUE)Fichiers créés :$(NC)"
	@ls -la $(BUILD_DIR)/
	@echo ""
	@echo "$(YELLOW)Pour utiliser le binder :$(NC)"
	@echo "  ./$(BINDER_BIN) <nom_sortie> <programme1> <programme2>"

# Cible pour le développement - reconstruction automatique
dev: build
	@echo "$(BLUE)Mode développement - surveillance des fichiers...$(NC)"
	@echo "$(YELLOW)Utilisez Ctrl+C pour arrêter$(NC)"
	@while inotifywait -e modify $(GO_FILES) 2>/dev/null; do \
		echo "$(YELLOW)Fichier modifié, recompilation...$(NC)"; \
		make build; \
	done

# Vérifier le formatage du code
fmt:
	@echo "$(BLUE)Vérification du formatage...$(NC)"
	@gofmt -l $(GO_FILES) | grep . && echo "$(RED)Fichiers mal formatés trouvés$(NC)" || echo "$(GREEN)✓ Code bien formaté$(NC)"

# Formater le code
fmt-fix:
	@echo "$(BLUE)Formatage du code...$(NC)"
	@gofmt -w $(GO_FILES)
	@echo "$(GREEN)✓ Code formaté$(NC)"

# Vérifier avec go vet
vet:
	@echo "$(BLUE)Vérification avec go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Aucun problème détecté$(NC)"

# Exécuter tous les contrôles de qualité
quality: fmt vet
	@echo "$(GREEN)✓ Contrôles de qualité terminés$(NC)"