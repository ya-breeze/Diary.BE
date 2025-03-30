{{ template "header.tpl" . }}

<main>
    <h2>Matchers</h2>
    {{ range .Matchers }}
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ .Name }}</h5>
            <a href="/web/matchers/edit?id={{.Id}}" class="btn btn-primary" tabindex="-1" role="button">Edit</a>
            <a href="/web/matchers/delete?id={{.Id}}" class="btn btn-danger" tabindex="-1" role="button">Delete</a>
        </div>
    </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}