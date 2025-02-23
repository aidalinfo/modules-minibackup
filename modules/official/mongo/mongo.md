# Module officiel MongoDB

---

Ce module permet de gérer les backups et la restauration des bases de données MongoDB.

## Configuration (Avec minibackup serveur)

Le fichier de configuration se trouve à l'emplacement `config/mongo.backups.yaml`.

### Exemple de configuration

```yaml
backups:
  mongo:
    type: mongo
    mongo: 
      host: "localhost"
      port: "27017"
      user: "root"
      password: "example"
      ssl: false
    path:
      local: "./backups"
      s3: "backup/mongo/mongo"
    retention:
      standard: 
        days: 14
    schedule:
      standard: "*/2 * * * *"
```


## Explication des clés de configuration

### backups

Contient la configuration de chaque backup.

### type

Spécifie le type de backup (ici, mongo).

### mongo

Contient la configuration de connexion MongoDB :

- `host` : L'hôte ou l'IP du serveur MongoDB.
- `port` : Le port du serveur MongoDB.
- `user` : Le nom d'utilisateur pour la connexion.
- `password` : Le mot de passe pour la connexion.
- `ssl` : Indique si SSL doit être utilisé (true/false).


### path

Définit les chemins de sauvegarde :

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

Pour lancer un backup MongoDB, utilisez la commande suivante :

```bash
module-mongo-mb backup [name] [args]
```


#### Arguments :

- `name` : Le nom de la configuration de backup (par exemple, `mongo` tel que défini dans le fichier YAML).
- `args` : Une chaîne JSON contenant les clés de configuration (au format YAML) que vous souhaitez appliquer ou surcharger.

**Exemple** :

```bash
module-mongo-mb backup mongo '{"mongo": {"host": "localhost", "port": "27017"}}'
```


### Restauration (restore)

Pour restaurer une base MongoDB, utilisez la commande suivante :

```bash
module-mongo-mb restore [name] [backupPath] [args]
```


#### Arguments :

- `name` : Le nom de la configuration de restauration (correspondant à une clé sous `backups` dans le fichier YAML).
- `backupPath` : Le chemin vers le dossier contenant le backup à restaurer.
- `args` : Une chaîne JSON avec les clés de configuration (format YAML) à appliquer pour la restauration.

**Exemple** :

```bash
module-mongo-mb restore mongo /path/to/backup '{"mongo": {"host": "localhost", "port": "27017"}}'
```


## Métadonnées du module

- **Version** : 0.1.0
- **Type** : mongo
- **Auteur** : Ninapepite
- **Description** : Module de sauvegarde MongoDB


## Options globales

```
  -h, --help   Affiche l'aide pour module-mongo-mb
```

**Note** : Les arguments `[args]` sont des chaînes JSON contenant des clés correspondant aux paramètres de votre fichier YAML de configuration. Assurez-vous que les clés utilisées correspondent bien à celles définies (par exemple, `mongo`, `path`, `retention`, etc.) afin de personnaliser correctement le comportement du backup ou de la restauration.

