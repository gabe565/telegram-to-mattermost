## telegram-to-mattermost



```
telegram-to-mattermost dir [flags]
```

### Options

```
      --allow-unknown-fields   Allows unknown fields in the Telegram JSON file. Unsupported Telegram fields may be silently ignored!
      --create-users           Adds users to Mattermost import (default true)
  -h, --help                   help for telegram-to-mattermost
      --max-text-length uint   Maximum post text length (default 4000)
      --no-attachments         Disables embedding of attachments
      --no-fix-webp            Disable fixing of WebP files (usually stickers) which will not load into Mattermost
  -o, --output string          Output filename (default "data.zip")
      --team-name string       Mattermost team name to import into
```

