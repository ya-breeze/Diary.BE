<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ block "title" . }}My Web App{{ end }}</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet"
        integrity="sha384-QWTKZyjpPEjISv5WaRU9OFeRpok6YctnYmDr5pNlyT2bRjXh0JMhjY6hW+ALEwIH" crossorigin="anonymous">
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap-icons@1.11.3/font/bootstrap-icons.min.css">
    <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/js/bootstrap.bundle.min.js"
        integrity="sha384-YvpcrYf0tY3lHB60NNkmXc5s9fDVZLESaAA55NDzOxhy9GkcIdslK1eN7N6jIeHz"
        crossorigin="anonymous"></script>
</head>

<body>
    <header>
        <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
            <div class="container-fluid">
                <a class="navbar-brand" href="/">GeekBudget Lite</a>
                <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav" aria-controls="navbarNav" aria-expanded="false" aria-label="Toggle navigation">
                    <span class="navbar-toggler-icon"></span>
                </button>
                <div class="collapse navbar-collapse" id="navbarNav">
                    <ul class="navbar-nav">
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "home"}}active{{end}}" href="/">Home</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "transactions"}}active{{end}}" href="/web/transactions">Transactions</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "accounts"}}active{{end}}" href="/web/accounts">Accounts</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "bank_importers"}}active{{end}}" href="/web/bank-importers">Importers</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "matchers"}}active{{end}}" href="/web/matchers">Matchers</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "unprocessed"}}active{{end}}" href="/web/unprocessed">Unprocessed</a>
                        </li>
                        <li class="nav-item">
                            <a class="nav-link {{if eq .CurrentPage "about"}}active{{end}}" href="/web/about">About</a>
                        </li>
                    </ul>
                    <form class="d-flex ms-auto">
                        <input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">
                        <button class="btn btn-outline-success" type="submit">Search</button>
                    </form>
                    <ul class="navbar-nav ms-3">
                        <li class="nav-item">
                            <a class="nav-link" href="/logout">Logout</a>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
