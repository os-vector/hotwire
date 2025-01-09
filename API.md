# API

- API definitions for Hotwire configuration, bot setup, and bot settings

## Hotwire dashboard

- Uptime, CPU usage, number of bots, etc
- Bots could maybe be pinnable? Could show bot status

## Hotwire configuration

- We need to configure:
  - Default Knowledge Graph/Intent Graph provider
  - Default Weather provider
  - Escape Pod or IP Address
    - Be able to provide IP address if needed
  - UI things (color, font)
- This probably doesn't need a websocket... we could just use a REST API

## Bot setup

- Need:
  - One "PAIR WITH VECTOR" section
  - The BLE code should be built in to Hotwire or the UI application
  - The SSH code should be built in there too
  - Section which shows a list of bots, allows you to remove them

## Bot settings (maybe "Bot Management")

- Need:
  - List of robots
    - Not just ESNs. Each robot should have it's own box in this section showing some information, like when the last command from the bot was and if the bot is active
    - Each robot could be given a name
  - A nice-looking, searchable settings interface for each bot
  - Vector Control section