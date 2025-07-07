package main

import (
	_ "embed"
	"log"
	"os"
	"os/exec"
)

//go:embed {{.ProgA}}
var progAbyte []byte

//go:embed {{.ProgB}}
var progBbyte []byte

// runEmbedded gère l'extraction et l'exécution d'un programme intégré.
func runEmbedded(name string, data []byte) {
	// Crée un fichier temporaire sécurisé dans le répertoire par défaut du système.[4, 5]
	tmpFile, err := os.CreateTemp("", name+"-exec-*")
	if err!= nil {
		log.Fatalf("Échec de la création du fichier temporaire pour %s: %v", name, err)
	}
	// Assure le nettoyage en supprimant le fichier temporaire à la fin de la fonction.[4, 5]
	defer os.Remove(tmpFile.Name())

	// Écrit les données binaires intégrées dans le fichier temporaire.
	if _, err := tmpFile.Write(data); err!= nil {
		log.Fatalf("Échec de l'écriture dans le fichier temporaire pour %s: %v", name, err)
	}
	// Ferme le fichier pour que le système d'exploitation puisse l'exécuter.
	if err := tmpFile.Close(); err!= nil {
		log.Fatalf("Échec de la fermeture du fichier temporaire pour %s: %v", name, err)
	}

	// Définit les permissions d'exécution, étape cruciale sur macOS et Linux.[6, 7, 8]
	if err := os.Chmod(tmpFile.Name(), 0755); err!= nil {
		log.Fatalf("Échec de la définition des permissions d'exécution pour %s: %v", name, err)
	}

	// Prépare la commande pour lancer l'exécutable.[9, 10]
	cmd := exec.Command(tmpFile.Name())
	// Connecte la sortie standard et l'erreur standard du sous-processus
	// à celles du processus principal pour que l'utilisateur voie la sortie.[9, 11, 12]
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Lance la commande et attend qu'elle se termine.
	// `Run` est bloquant, ce qui garantit une exécution séquentielle.[9, 13]
	if err := cmd.Run(); err!= nil {
		// Affiche une erreur non fatale si un des programmes échoue.
		log.Printf("Erreur lors de l'exécution de %s: %v", name, err)
	}
}

func main() {
	runEmbedded("programmeA", progAbyte)
	runEmbedded("programmeB", progBbyte)
}