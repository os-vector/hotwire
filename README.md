# Hotwire

- A soon-to-be rewrite of [wire-pod](https://github.com/kercre123/wire-pod).
- See the other MDs for implementation details.

## Things to implement

- [ ] Token server
- [ ] Jdocs server
- [ ] Chipper server
- [ ] STT interface
  - Should implement:
    - `Name string`
    - `Init() error`
    - `Process(data []int16, data []byte) (string, error)`
    - `MultiLanguageSupport bool`
       - If true, define SupportedLanguages. That could be a struct containing language, link to model, and accuracy notes
- [ ] Better VAD
  - Use same library, just handle more edge cases
- [ ] Better voice filtering
  - wire-pod's voice filtering seems to slow things down...
- [ ] Websocketed API
  - [ ] Configuration
  - [ ] Bot settings
  - Maybe just use websockets for robot control
- [ ] UI made with a framework
  - React? Flutter?
  - pretty
- [ ] Inbuilt BLE with go-ble
  - [ ] Detect a dev OS. If dev, download logs for SSH key (after trying Anki ssh_root_key) and implement correct vic-cloud and server config
  - [ ] Handle OTA updates. Consider uploading the OTA to a GitHub release and including an HTTPS->HTTP proxy in Hotwire 
- [ ] Make connchecks work, will require gRPC-gateway
- [ ] Handle IP changes (both robot and Hotwire server)
- [ ] Use combo of multicast and mDNS?
- [ ] Easy-to-use "get GUID" API endpoint
- [ ] Configurable intent utterances
  - [ ] Involves a list of commands and descriptions
- [ ] Knowledge graph
  - [ ] Better StreamingKGSim function which actually handles errors
  - [ ] More LLM commands
  - [ ] Make conversations happen more
- [ ] Weather
- [ ] More configurable weather and knowledge graph
  - Use interfaces to swap in APIs
  - [ ] Gemini?
- [ ] Test API keys on interface
- [ ] Interface should redirect user to API provider login and get API keys via a more correct way
- [ ] Security (username/password)?
- [ ] Have Bot Settings be its own project which can be directed towards a different server?
- [ ] Allow importing data from wire-pod

## Releases

- [ ] Windows
- [ ] macOS ARM and x86
- [ ] Android
  - [ ] a better app which creates a foreground service
- [ ] Debian/Ubuntu
- [ ] Docker (w/ storage and mDNS instructions)
- [ ] iOS??
- [ ] AUR

## Progress

- Copied in token, jdocs, and chipper servers from cavalier
- Very broken. This only contains ideas and experiments

## Repo owners

- [kercre123/Wire](https://github.com/kercre123)
