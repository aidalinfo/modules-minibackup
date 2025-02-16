#!/bin/bash

# Définir le chemin de base
BASE_DIR="$(dirname "$(realpath "$0")")"
MODULES_DIR="$BASE_DIR/modules"

# Vérifier si jq est installé
if ! command -v jq &>/dev/null; then
    echo "❌ Erreur : 'jq' n'est pas installé. Installez-le avec 'sudo apt install jq' (ou équivalent)."
    exit 1
fi

# Initialisation de la structure JSON
INDEX_JSON='{
    "official": {},
    "community": {},
    "collections": {}
}'

# Fonction pour extraire une clé YAML (compatibilité Bash sans yq)
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

# Fonction pour extraire toutes les clés sous metadata en JSON
extract_metadata() {
    local file=$1
    awk '/^ *metadata:/,/^[^ ]/{if ($1 !~ /metadata:/) print}' "$file" |
        sed -E 's/^[ ]+([^:]+):[ ]*(.*)$/"\1": "\2",/' |
        tr -d '\n' | sed 's/,$//'
}

# Fonction pour parcourir les modules et remplir le JSON
process_modules() {
    local category=$1
    local category_path="$MODULES_DIR/$category"

    if [[ ! -d "$category_path" ]]; then
        return
    fi

    for module_dir in "$category_path"/*; do
        if [[ -d "$module_dir" ]]; then
            local yaml_file="$module_dir/module.yaml"

            if [[ -f "$yaml_file" ]]; then
                local name version type metadata_json
                name=$(extract_module_name "$yaml_file")
                version=$(extract_yaml_value "version" "$yaml_file")
                type=$(extract_yaml_value "type" "$yaml_file")
                metadata_json=$(extract_metadata "$yaml_file")

                local md_file="$module_dir/$name.md"
                if [[ -f "$md_file" && -n "$name" && -n "$version" && -n "$type" ]]; then
                    local path="modules/$category/$(basename "$module_dir")/$(basename "$md_file")"
                    
                    if [[ -n "$metadata_json" ]]; then
                        INDEX_JSON=$(echo "$INDEX_JSON" | jq --arg category "$category" \
                                                              --arg name "$name" \
                                                              --arg version "$version" \
                                                              --arg type "$type" \
                                                              --arg path "$path" \
                                                              --argjson metadata "{$metadata_json}" \
                                                              '.[$category][$name] = { "version": $version, "type": $type, "path": $path, "metadata": $metadata }')
                    else
                        INDEX_JSON=$(echo "$INDEX_JSON" | jq --arg category "$category" \
                                                              --arg name "$name" \
                                                              --arg version "$version" \
                                                              --arg type "$type" \
                                                              --arg path "$path" \
                                                              '.[$category][$name] = { "version": $version, "type": $type, "path": $path, "metadata": {} }')
                    fi
                fi
            fi
        fi
    done
}

# Parcourir les catégories
process_modules "official"
process_modules "community"
process_modules "collections"

# Enregistrer dans .index.json
echo "$INDEX_JSON" | jq '.' > "$BASE_DIR/.index.json"

echo "✅ Fichier .index.json généré avec succès !"
