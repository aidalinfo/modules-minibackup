# Module Filesystem pour MiniBackup

---

Ce module permet de gérer les backups et la restauration des systèmes de fichiers.

## Configuration

Le fichier de configuration se trouve à l'emplacement `fs.backups.yaml`.

### Exemple de configuration

```yaml
backups:
  filesystem:
    type: fs
    fs:
      paths:
        - "./datatest/fs"
    path:
      local: "./backups"
      s3: "minio-data"
    retention:
      standard:
        days: 14
    schedule:
      standard: "*/59 * * * *"
```


## Explication des clés de configuration

### backups

Contient la configuration de chaque backup.

### type

Spécifie le type de backup (ici, fs pour filesystem).

### fs

Contient la configuration du backup filesystem :

- `paths` : Liste des chemins à sauvegarder.


### path

Définit les chemins de sauvegarde :

- `local` : Le répertoire local pour stocker les backups.
- `s3` : Le chemin de destination sur S3 ou un système compatible.


### retention

Définit la politique de rétention des backups.

- `standard` : Politique de rétention par défaut.
- `days` : Nombre de jours pendant lesquels conserver les backups.


### schedule

Spécifie la planification des backups (format cron).

## Usage manuel

### Sauvegarde (backup)

Pour lancer un backup filesystem, utilisez la commande suivante :

```bash
module-fs-mb backup filesystem '{"fs": {"paths": ["./datatest/fs"]}}'
```


#### Arguments :

- `name` : Le nom de la configuration de backup (par exemple, `filesystem` tel que défini dans le fichier YAML).
- `args` : Une chaîne JSON contenant les clés de configuration que vous souhaitez appliquer ou surcharger.


### Restauration (restore)

Pour restaurer un système de fichiers, utilisez la commande suivante :

```bash
module-fs-mb restore filesystem /path/to/backup '{"fs": {"paths": ["./datatest/fs"]}}'
```


#### Arguments :

- `name` : Le nom de la configuration de restauration (correspondant à une clé sous `backups` dans le fichier YAML).
- `backupPath` : Le chemin vers le dossier contenant le backup à restaurer.
- `args` : Une chaîne JSON avec les clés de configuration à appliquer pour la restauration.


## Métadonnées du module

- **Version** : 0.1.0
- **Auteur** : Ninapepite
- **Description** : Module de sauvegarde Filesystem
