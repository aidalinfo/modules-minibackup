# Official MySQL Module

Ce module permet de gérer les backups et la restauration des bases de données MySQL.

## Configuration (Avec minibackup serveur)

Le fichier de configuration se trouve à l'emplacement `config/mysql.backups.yaml`.

### Exemple de configuration

```yaml
backups:
  sqlserver-01:
    type: mysql
    mysql:
      all: true
      host: "localhost"
      port: "3306"
      user: "root"
      password: "example"
      ssl: "false"
    path:
      local: "./backups"
      s3: "backup/glpi-dev/mysql"
    retention:
      standard:
        days: 20
    schedule:
      standard: "*/1 * * * *"
```

## Explication des clés de configuration

### backups
Contient la configuration de chaque backup.

### type
Spécifie le type de backup (ici, mysql).

### mysql
Contient la configuration de connexion MySQL :

- `all` : Si `true`, toutes les bases seront sauvegardées.
- `host` : L'hôte ou l'IP du serveur MySQL.
- `port` : Le port du serveur MySQL.
- `user` : Le nom d'utilisateur pour la connexion.
- `password` : Le mot de passe pour la connexion.
- `ssl` : Indique si SSL doit être utilisé.
- `databases` : Si `all` est `false`, liste des bases spécifiques à sauvegarder.

### path
Définit les chemins de sauvegarde :

- `local` : Le répertoire local pour stocker les backups.
- `s3` : Le chemin de destination sur S3.

### retention
Définit la politique de rétention des backups.

- `standard` : Politique de rétention par défaut.
- `days` : Nombre de jours pendant lesquels conserver les backups.

### schedule
Spécifie la planification des backups (format cron).

## Usage manuel (sans minibackup serveur)

### Sauvegarde (backup)
Pour lancer un backup MySQL, utilisez la commande suivante :

```bash
minibackup module-mysql-mb backup sqlserver-01 '{"mysql": {"all": true}}'
```

#### Arguments :

- `name` : Le nom de la configuration de backup (par exemple, `sqlserver-01` tel que défini dans le fichier YAML).
- `args` : Une chaîne JSON contenant les clés de configuration (au format YAML) que vous souhaitez appliquer ou surcharger.

**Exemple** : `{ "mysql": { "all": true } }`

### Restauration (restore)
Pour restaurer une base MySQL, utilisez la commande suivante :

```bash
minibackup module-mysql-mb restore sqlserver-01 /path/to/backup '{"mysql": {"all": true}}'
```

#### Arguments :

- `name` : Le nom de la configuration de restauration (correspondant à une clé sous `backups` dans le fichier YAML).
- `backupPath` : Le chemin vers le dossier contenant le backup à restaurer.
- `args` : Une chaîne JSON avec les clés de configuration (format YAML) à appliquer pour la restauration.


Cette commande créera un répertoire `docs` contenant des fichiers Markdown pour chaque commande (par exemple, `module-mysql-mb_backup.md` et `module-mysql-mb_restore.md`).

**Note** : Les arguments `{args}` sont des chaînes JSON contenant des clés correspondant aux paramètres de votre fichier YAML de configuration. Assurez-vous que les clés utilisées correspondent bien à celles définies (par exemple, `mysql`, `path`, `retention`, etc.) afin de personnaliser correctement le comportement du backup ou de la restauration.

