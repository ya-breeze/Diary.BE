{{ template "header.tpl" . }}

<main>
    <h2>Upload result</h2>

    {{ with .LastImport }}
        <h3>Date: {{ formatTime .Date "2006-01-02" }}</h3>
        <h3>Status: {{ .Status }}</h3>

        <h3>Description</h3>
        <p>{{ .Description }}</p>
    {{ end }}
</main>

{{ template "footer.tpl" . }}