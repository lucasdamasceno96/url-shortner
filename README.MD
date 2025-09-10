# 🔗 URL Shortener API in Go  

This project is a functional **URL shortening API** developed in **Go** 🐹.  
It takes a long URL and generates a unique, shorter alias 🔑 that redirects to the original URL 🌍.  

The main goal is not only to build the service but to do so using the **best software engineering practices** 🛠️,  
creating a **clean, maintainable, and testable** codebase.  

---

## ✨ Features  
- 🔒 **URL Shortening**: Converts long URLs into a compact format.  
- ↪️ **URL Redirection**: Redirects short codes to their original URLs.  
- 💾 **Persistent Storage**: Uses SQLite to store URL mappings.  
- 🌐 **RESTful API**: Provides simple and clear API endpoints.  
- 🧩 **Layered Architecture**: Ensures separation of concerns and modularity.  

---

## 🏗️ Tech Stack & Architecture  
- 🐹 **Language**: Go (Golang)  
- 📡 **API Framework**: gorilla/mux for HTTP routing  
- 🗄️ **Database**: SQLite 3 (lightweight, file-based)  
- 🔌 **Driver**: mattn/go-sqlite3  

The app follows **Clean Architecture** 🧼, with 4 distinct layers:  
- 📦 **Domain** → Core data structures (e.g., `ShortURL`). Independent from other layers.  
- 🗄️ **Repository** → Handles DB communication.  
- ⚙️ **Service** → Business logic (e.g., generate unique short codes, validate input).  
- 🌍 **Handler (HTTP)** → Handles HTTP requests/responses.  

---

## 📂 Project Structure  
