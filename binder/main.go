package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

// Chemin vers le modèle du stub de chargeur.
const loaderStubPath = "loader/stub.go.tpl"

func main() {
	// 1. Vérifier les arguments de la ligne de commande.
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run. <nom_sortie> <chemin_programme1> <chemin_programme2>")
		os.Exit(1)
	}
	outputName := os.Args[1]
	progAPath := os.Args[2]
	progBPath := os.Args[3]

	// 2. Créer un répertoire de construction temporaire.[4, 5]
	buildDir, err := os.MkdirTemp("", "binder-build-*")
	if err != nil {
		log.Fatalf("Échec de la création du répertoire de build temporaire: %v", err)
	}
	defer os.RemoveAll(buildDir)
	log.Printf("Utilisation du répertoire de build temporaire : %s", buildDir)

	// 3. Lire les exécutables d'entrée et les écrire dans le répertoire de build.
	payloadA, err := os.ReadFile(progAPath)
	if err != nil {
		log.Fatalf("Échec de la lecture du programme A : %v", err)
	}
	err = os.WriteFile(filepath.Join(buildDir, "progA.bin"), payloadA, 0644)
	if err != nil {
		log.Fatalf("Échec de l'écriture du payload du programme A : %v", err)
	}

	payloadB, err := os.ReadFile(progBPath)
	if err != nil {
		log.Fatalf("Échec de la lecture du programme B : %v", err)
	}
	err = os.WriteFile(filepath.Join(buildDir, "progB.bin"), payloadB, 0644)
	if err != nil {
		log.Fatalf("Échec de l'écriture du payload du programme B : %v", err)
	}

	// 4. Générer le fichier main.go pour l'exécutable final à partir du modèle.[14, 15, 16]
	tpl, err := template.ParseFiles(loaderStubPath)
	if err != nil {
		log.Fatalf("Échec de l'analyse du modèle de stub : %v", err)
	}

	generatedGoFile, err := os.Create(filepath.Join(buildDir, "main.go"))
	if err != nil {
		log.Fatalf("Échec de la création du fichier main.go généré : %v", err)
	}

	templateData := struct {
		ProgA string
		ProgB string
	}{
		ProgA: "progA.bin",
		ProgB: "progB.bin",
	}

	err = tpl.Execute(generatedGoFile, templateData)
	if err != nil {
		log.Fatalf("Échec de l'exécution du modèle : %v", err)
	}
	generatedGoFile.Close()

	// Copier go.mod et go.sum dans le dossier temporaire
	copyFile := func(src, dst string) error {
		input, err := os.ReadFile(src)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, input, 0644)
	}

	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Impossible d'obtenir le répertoire courant : %v", err)
	}
	err = copyFile(filepath.Join(projectRoot, "go.mod"), filepath.Join(buildDir, "go.mod"))
	if err != nil {
		log.Fatalf("Impossible de copier go.mod : %v", err)
	}
	// Copier go.sum si présent
	if _, err := os.Stat(filepath.Join(projectRoot, "go.sum")); err == nil {
		_ = copyFile(filepath.Join(projectRoot, "go.sum"), filepath.Join(buildDir, "go.sum"))
	}

	// 5. Compiler le code généré pour créer l'exécutable final.[10, 17, 18]
	log.Println("Compilation de l'exécutable final...")
	outputPath, err := filepath.Abs(outputName)
	if err != nil {
		log.Fatalf("Échec de l'obtention du chemin absolu pour la sortie : %v", err)
	}

	cmd := exec.Command("go", "build", "-o", outputPath, ".")
	cmd.Dir = buildDir // Exécuter la compilation dans notre répertoire temporaire.
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("La compilation Go a échoué : %v", err)
	}

	log.Printf("Exécutable créé avec succès : %s", outputPath)
}
