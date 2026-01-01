# Groupie Tracker

**Groupie Tracker** is a web application written in Go that fetches data from a public API and displays information about music artists and bands. The application shows details such as:

* **Name**
* **First Album**
* **Formation Date**
* **Members**
* **Concerts** (locations and dates)

---

## Features

* The home page displays all artists as cards with images and names.
* Clicking on an artist card opens a detailed page with all available information.
* Concert locations and dates are shown for each artist.
* Error handling for invalid requests or missing data.
* Uses only **Go standard packages**, without any external libraries.

---

## How to Run

1. Clone the repository:

   ```bash
   git clone <repository-url>
   ```
2. Make sure Go is installed on your system.
3. Run the server:

   ```bash
   go run .
   ```
4. Open your browser and go to:

   ```
   http://localhost:8080
   ```

---

## Project Structure

```
.
├── checkdata.txt
├── go.mod
├── handlers/
│   ├── artist.go
│   ├── date.go
│   ├── Error.go
│   ├── fetchaartist.go
│   ├── fetchlocaltion.go
│   ├── home.go
│   ├── Relations.go
│   ├── statichandler.go
│   └── sub-func.go
├── main.go
├── README.md
├── static/
│   ├── artist.css
│   ├── error.css
│   └── index.css
└── templates/
    ├── artist.html
    ├── err.html
    └── index.html
```

* `main.go` – the entry point that sets up the server and routes.
* `handlers/` – contains Go files for HTTP handlers, API fetching, and page rendering.
* `templates/` – HTML templates for the home page and artist detail pages.
* `static/` – CSS files for styling the pages.

---

## API

The application uses the [Groupie Tracker API](https://groupietrackers.herokuapp.com/api) to fetch all artist, location, date, and relation data.

---

## License

This project is for educational purposes only.

---

## Authors

- yzarhoun  - [GitHub](https://github.com/quadS01) - [Gitea](https://learn.zone01oujda.ma/git/yzarhoun)
- halhyane  - [GitHub](https://github.com/Houssam-Alhyane) - [Gitea](https://learn.zone01oujda.ma/git/halhyane)