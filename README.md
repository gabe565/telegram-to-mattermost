# Telegram to Mattermost

A command-line tool that converts Telegram exports to Mattermost data imports.

## Usage

1. Install the [Telegram Desktop app](https://desktop.telegram.org).
2. Export chat history:
   1. Select a chat
   2. Open the overflow menu in the top-right
   3. Export chat history
   4. Choose the options you desire, but ensure "Format" is set to "Machine-readable JSON"
3. Install the tool.
   ```shell
   go install github.com/gabe565/telegram-to-mattermost@latest
   ```
4. Run the tool with the data export directory as a param, for example:
   ```shell
   telegram-to-mattermost ~/Downloads/Telegram\ Desktop/ChatExport_2024-08-01
   ```
   - See [docs](./docs/telegram-to-mattermost.md) for command-line flags.
5. Fill in the prompts to set up user mappings.
6. The tool should finish successfully, and a new file `data.zip` will be available to import into Mattermost.
7. Follow Mattermost's [bulk loading data](https://docs.mattermost.com/onboard/bulk-loading-data.html#bulk-load-data) guide to import this file.
