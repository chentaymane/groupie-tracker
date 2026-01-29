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
├── fetchers
│   ├── fetchArtist.go
│   ├── fetchDates.go
│   ├── fetchLocaltion.go
│   ├── fetchRelations.go
│   └── sub-func.go
├── go.mod
├── handlers
│   ├── artistHandler.go
│   ├── errorHandler.go
│   ├── homeHandler.go
│   └── staticHandler.go
├── main.go
├── README.md
├── static
│   └── style.css
└── templates
    ├── artist.html
    ├── err.html
    └── index.html
```
---

## API

The application uses the [Groupie Tracker API](https://groupietrackers.herokuapp.com/api) to fetch all artist, location, date, and relation data.

---

## License

This project is for educational purposes only.

---

## Authors

- mbelhouss  - [GitHub](https://github.com/DissonantVoid) - [Gitea](https://learn.zone01oujda.ma/git/mbelhouss)
- achent  - [GitHub](https://github.com/chentaymane) - [Gitea](https://learn.zone01oujda.ma/git/achent)