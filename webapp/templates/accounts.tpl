{{ template "header.tpl" . }}

<main>
    <h2>Accounts</h2>

    <a class="link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
        href="/web/accounts/edit">
        Add account
    </a>

    {{ range .Accounts }}
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ .Name }}</h5>
            <p class="card-text">{{ .Type }}</p>
            <p class="card-text">{{ .Description }}</p>
            <a class="link-offset-2 link-offset-3-hover link-underline link-underline-opacity-0 link-underline-opacity-75-hover"
                href="/web/accounts/edit?id={{.Id}}">
                Edit account
            </a>
            <form action="/web/accounts" method="DELETE">
                <input type="hidden" name="id" value="{{ .Id }}">
                <button type="submit" class="btn btn-primary">Delete</button>
            </form>
        </div>
    </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}