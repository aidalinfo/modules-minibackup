Example usign the official mysql module

## Warning
Please actually not use glacier mode for mysql backups, it is not supported by the official mysql module.

```yaml
backups:
  sqlserver-01:
    type: mysql
    mysql:
      allDatabases: true
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