<br />
<div align="center">

<h3 align="center">Sengen</h3>

  <p align="center">
    Sengen is an AI anki app that uses Gemini and Google TTS to generate sentences, translations, definitions, and add them to your anki decks automatically.
    <br />
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#run-locally">Run Locally</a></li>
        <li><a href="#run-with-docker">Run with Docker</a></li>
      </ul>
    </li>
    <li><a href="#api">API</a></li>
    <li><a href="#under-the-hood">Under the Hood</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

This is a backend gRPC service built to power [Sengen](https://github.com/dafraer/sentence-gen-grpc-client)  language learning app. It serves as the AI engine behind the tool — handling sentence generation, translation, and word definitions.

Given a word and a target language, Sengen can:

- **Generate a contextual example sentence** — along with its translation into another language, so learners can understand meaning from context
- **Translate a word or phrase** — between any two languages, with an optional hint to disambiguate meaning
- **Define a word** — producing a clean, monolingual definition without restating the word itself

All three operations optionally return **audio pronunciation** synthesized via Google's Chirp3-HD text-to-speech model.

Sengen also enforces a configurable **daily spending quota**, tracking Gemini token usage and TTS character counts in Firestore to keep API costs under control.



<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- **Go 1.25+**
- **Google Cloud project** with Firestore and Text-to-Speech APIs enabled
- **Google Application Default Credentials** configured — follow [this guide](https://cloud.google.com/docs/authentication/application-default-credentials)
- **Gemini API key** — get one at [Google AI Studio](https://aistudio.google.com/apikey)

> **Note:** The Gemini and Google Cloud APIs are not free. Make sure you understand the pricing before running in production.

### Run Locally

#### 1. Clone the repository

```sh
git clone https://github.com/dafraer/sentence-gen-grpc-server.git
cd sentence-gen-grpc-server
```

#### 2. Configure environment variables

Copy the example env file and fill in your values:

```sh
cp .env_example .env
```

```env
GEMINI_API_KEY=<your-api-key>
PROJECT_ID=<your-gcp-project-id>
ADDRESS=localhost:50051
GEMINI_MODEL=gemini-2.5-pro-preview
DAILY_QUOTA=5000000        # Daily spending cap in micro USD
GEMINI_INPUT_PRICE=2       # Price per input token in micro USD
GEMINI_OUTPUT_PRICE=12     # Price per output token in micro USD
```

#### 3. Start the server

```sh
make run
```

The gRPC server will be listening on the address specified in `ADDRESS` (default: `localhost:50051`).

### Run with Docker

Build and run the container locally:

```sh
docker build -t sengen .
docker run --env-file .env -p 50051:50051 sengen
```

Or pull the published image:

```sh
docker pull dafraer/sentence-gen-grpc-server:<version>
docker run --env-file .env -p 50051:50051 dafraer/sentence-gen-grpc-server:<version>
```



<!-- API -->
## API

Sengen exposes a single gRPC service defined in [`proto/sentence-gen.proto`](proto/sentence-gen.proto):

```protobuf
service SentenceGen {
  rpc GenerateSentence(GenerateSentenceRequest) returns (GenerateSentenceResponse);
  rpc Translate(TranslateRequest) returns (TranslateResponse);
  rpc GenerateDefinition(GenerateDefinitionRequest) returns (GenerateDefinitionResponse);
}
```

### `GenerateSentence`

Generates a contextual example sentence for a word and its translation.

| Field | Type | Description |
|---|---|---|
| `word_language` | string | Language of the word (e.g. `"Japanese"`) |
| `translation_language` | string | Language for the translation (e.g. `"English"`) |
| `word` | string | The word to use in the sentence |
| `translation_hint` | string | Optional hint to disambiguate meaning |
| `include_audio` | bool | Whether to include audio of the sentence |
| `voice_gender` | Gender | `GENDER_FEMALE` or `GENDER_MALE` |

Returns `original_sentence`, `translated_sentence`, and optionally `audio` (WAV bytes).

### `Translate`

Translates a word or phrase between two languages.

| Field | Type | Description |
|---|---|---|
| `from_language` | string | Source language |
| `to_language` | string | Target language |
| `word` | string | Word or phrase to translate |
| `translation_hint` | string | Optional disambiguation hint |
| `include_audio` | bool | Whether to include audio of the source word |
| `voice_gender` | Gender | `GENDER_FEMALE` or `GENDER_MALE` |

Returns `translation` and optionally `audio` (WAV bytes).

### `GenerateDefinition`

Generates a monolingual definition of a word.

| Field | Type | Description |
|---|---|---|
| `language` | string | Language of the word |
| `word` | string | Word or term to define |
| `definition_hint` | string | Optional hint to guide the definition |
| `include_audio` | bool | Whether to include audio of the word |
| `voice_gender` | Gender | `GENDER_FEMALE` or `GENDER_MALE` |

Returns `definition` and optionally `audio` (WAV bytes).

To regenerate the protobuf bindings after modifying the `.proto` file:

```sh
make generate
```



<!-- UNDER THE HOOD -->
## Under the Hood

Here's a breakdown of the tech powering Sengen:

- **Core Logic**
  Written in **Go**, exposing a clean gRPC interface with a unary interceptor for daily quota enforcement.

- **Language Model**
  Uses [**Gemini**](https://gemini.google.com/) via the Google GenAI SDK with **structured JSON output** 

- **Text-to-Speech**
  Audio is generated using the [**Google Cloud Text-to-Speech API**](https://cloud.google.com/text-to-speech) with the **Chirp3-HD** neural voice model, producing high-quality WAV audio. Voice selection is dynamic — the server queries available voices for the requested language and gender at runtime, and gracefully skips audio if no matching voice exists.

- **Database**
  [**Google Firestore**](https://firebase.google.com/docs/firestore) is used to persist daily API spending, enabling the quota limiter to track Gemini token usage and TTS character counts across requests.

- **Quota Limiter**
  A gRPC unary interceptor checks daily spending against a configurable `DAILY_QUOTA` before every request. Costs are calculated in **micro USD** per token (Gemini) and per character (TTS), and accumulated atomically in Firestore.

- **Logging**
  Structured logging via [**go.uber.org/zap**](https://pkg.go.dev/go.uber.org/zap) throughout all layers.

- **Deployment**
  Containerized with a multi-stage **Docker** build (Alpine-based), producing a minimal image that exposes port `50051`. Published to Docker Hub as [`dafraer/sentence-gen-grpc-server`](https://hub.docker.com/r/dafraer/sentence-gen-grpc-server).



<!-- CONTACT -->
## Contact

Kamil Nuriev — [Telegram](https://t.me/dafraer) — kdnuriev@gmail.com
