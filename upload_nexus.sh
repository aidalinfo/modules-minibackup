#!/bin/bash

# Configuration de Nexus Repository
NEXUS_URL="https://pkg.aidalinfo.fr/repository/minibackup-modules"

# Définition du répertoire des modules
BASE_DIR="$(dirname "$(realpath "$0")")"
MODULES_DIR="$BASE_DIR/modules"
VERSIONS_FILE="$BASE_DIR/.versions.json"

# Vérification des dépendances
if ! command -v jq &>/dev/null; then
    echo "❌ 'jq' n'est pas installé. Installez-le avec 'sudo apt install jq'"
    exit 1
fi

# Initialisation du fichier .versions.json s'il n'existe pas
if [[ ! -f "$VERSIONS_FILE" ]]; then
    echo '{}' > "$VERSIONS_FILE"
fi

# Fonction pour extraire une clé YAML
extract_yaml_value() {
    local key=$1
    local file=$2
    grep -E "^ *$key:" "$file" | awk -F': ' '{print $2}' | tr -d ' '
}

# Fonction pour récupérer le nom du module
extract_module_name() {
    local file=$1
    awk -F ':' '/^[a-zA-Z0-9_-]+:/ {print $1; exit}' "$file" | tr -d ' '
}

# Fonction pour récupérer la version actuelle depuis .versions.json
get_local_version() {
    local name=$1
    jq -r --arg name "$name" '.[$name].version // "0.0.0"' "$VERSIONS_FILE"
}

version_greater() {
    # Retourne 1 (false) si les versions sont identiques
    if [ "$1" = "$2" ]; then
        return 1
    fi
    # Sinon, compare en utilisant sort -V
    if [ "$(printf '%s\n%s' "$1" "$2" | sort -V | tail -n 1)" = "$1" ]; then
        return 0
    else
        return 1
    fi
}

# Fonction principale : Build, ZIP et Upload
process_module() {
    local category=$1
    local module_path=$2
    local yaml_file="$module_path/module.yaml"

    if [[ ! -f "$yaml_file" ]]; then
        echo "⚠️ Aucun fichier module.yaml trouvé dans $module_path"
        return
    fi

    local name version bin_file zip_file
    name=$(extract_module_name "$yaml_file")  # Nom du module basé sur la première clé YAML
    version=$(extract_yaml_value "version" "$yaml_file")
    bin_file=$(extract_yaml_value "bin" "$yaml_file")

    if [[ -z "$name" || -z "$version" || -z "$bin_file" ]]; then
        echo "⚠️ Fichier module.yaml incomplet pour $module_path"
        return
    fi

    zip_file="$BASE_DIR/${name}.zip"

    # Récupérer la version locale stockée dans .versions.json
    local_version=$(get_local_version "$name")

    # Vérification si la version doit être mise à jour
    if [[ -n "$local_version" ]] && ! version_greater "$version" "$local_version"; then
        echo "✅ La version $version est inférieure ou égale à $local_version."
        return
    fi

    # Vérifier si c'est un module Go et compiler si nécessaire
    if [[ -f "$module_path/go.mod" ]]; then
        echo "🛠️ Détection d'un projet Go, compilation..."

        pushd "$module_path" > /dev/null
        go mod tidy
        go build -o "$bin_file" "cmd/main.go"
        popd > /dev/null

        if [[ $? -ne 0 ]]; then
            echo "❌ Échec de la compilation Go pour $name"
            return
        fi
    fi

    # Vérification du fichier binaire
    if [[ ! -f "$module_path/$bin_file" ]]; then
        echo "❌ Fichier binaire introuvable après compilation : $module_path/$bin_file"
        return
    fi

    # Création du fichier ZIP avec options pour garantir un SHA256 constant
    echo "📦 Création de l'archive ZIP $zip_file..."
    zip -X -j "$zip_file" "$module_path/$bin_file" "$yaml_file"

    if [[ ! -f "$zip_file" ]]; then
        echo "❌ Échec de la création de l'archive ZIP : $zip_file"
        return
    fi

    # Supprimer l'artefact existant sur Nexus pour écraser la version antérieure
    echo "⏳ Suppression de l'ancienne version sur Nexus (si existante)..."
    curl -u "$NEXUS_USERNAME:$NEXUS_PASSWORD" -X DELETE "$NEXUS_URL/$category/$name.zip" >/dev/null 2>&1

    # Upload vers Nexus (ZIP)
    echo "🚀 Upload du fichier sur Nexus..."
    curl -u "$NEXUS_USERNAME:$NEXUS_PASSWORD" \
        --upload-file "$zip_file" \
        "$NEXUS_URL/$category/$name.zip"

    if [[ $? -eq 0 ]]; then
        echo "✅ Upload réussi pour $name ($version)"
        
        # Calculer le SHA256 du fichier ZIP
        sha=$(sha256sum "$zip_file" | awk '{print $1}')
        
        # Mettre à jour le fichier .versions.json avec la version et le SHA
        jq --arg name "$name" --arg version "$version" --arg sha "$sha" \
            '.[$name] = { "version": $version, "sha": $sha }' "$VERSIONS_FILE" > "${VERSIONS_FILE}.tmp"

        mv "${VERSIONS_FILE}.tmp" "$VERSIONS_FILE"
        echo "✅ Fichier .versions.json mis à jour avec succès !"
        rm -f "$zip_file"  # Supprimer le fichier ZIP après l'upload
    else
        echo "❌ Échec de l'upload pour $name"
    fi
}

# Parcours des modules dans les catégories official, community, collections
for category in "official" "community" "collections"; do
    category_path="$MODULES_DIR/$category"

    if [[ -d "$category_path" ]]; then
        for module in "$category_path"/*; do
            if [[ -d "$module" ]]; then
                process_module "$category" "$module"
            fi
        done
    fi
done

echo "✅ Processus d'upload des modules terminé !"
