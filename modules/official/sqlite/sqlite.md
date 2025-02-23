# Module SQLite pour MiniBackup

---

Ce module permet de gérer les backups et la restauration des bases de données SQLite.

## Configuration

Le fichier de configuration se trouve à l'emplacement `config/sqlite.backups.yaml`.

### Exemple de configuration

```yaml
backups:
  sqlite-new:
    type: sqlite
    sqlite:
      paths: 
        - "./datatest/sqlite.db"
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

Spécifie le type de backup (ici, sqlite).

### sqlite

Contient la configuration du backup SQLite :

- `paths` : Liste des chemins des fichiers de base de données SQLite à sauvegarder.


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

Pour lancer un backup SQLite, utilisez la commande suivante :

```bash
module-sqlite-mb backup sqlite-new '{"sqlite": {"paths": ["./datatest/sqlite.db"]}}'
```


#### Arguments :

- `name` : Le nom de la configuration de backup (par exemple, `sqlite-new` tel que défini dans le fichier YAML).
- `args` : Une chaîne JSON contenant les clés de configuration que vous souhaitez appliquer ou surcharger.


### Restauration (restore)

Pour restaurer une base de données SQLite, utilisez la commande suivante :

```bash
module-sqlite-mb restore sqlite-new /path/to/backup '{"sqlite": {"paths": ["./datatest/sqlite.db"]}}'
```


#### Arguments :

- `name` : Le nom de la configuration de restauration (correspondant à une clé sous `backups` dans le fichier YAML).
- `backupPath` : Le chemin vers le dossier contenant le backup à restaurer.
- `args` : Une chaîne JSON avec les clés de configuration à appliquer pour la restauration.


## Métadonnées du module

- **Version** : 0.1.0
- **Auteur** : Ninapepite
- **Description** : Module de sauvegarde SQLite
