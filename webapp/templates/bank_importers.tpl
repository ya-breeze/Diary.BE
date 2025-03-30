{{ template "header.tpl" . }}

<main>
    <h2>Bank Importers</h2>
    {{ range .BankImporters }}
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ .Name }}</h5>
            <p class="card-text">{{ .Description }}</p>
            <p><strong>Last successful import:</strong> {{ formatTime .LastSuccessfulImport "2006-01-02 15:04" }}</p>
            {{ if ne .Type "fio" }}
                <form action="/web/bank-importers/upload" method="post" enctype="multipart/form-data">
                    <input type="hidden" name="id" value="{{ .Id }}">
                    <label for="file">Choose file to upload</label>
                    <input type="file" id="file" name="file">
                    <button class="btn btn-primary" type="submit">Upload</button>
                </form>
            {{ else }}
                {{ if eq .FetchAll true }}
                    Fetching of all transactions is in progress
                {{ else }}
                    <a href="?id={{.Id }}&fetchAll=true" class="btn btn-primary" >Fetch all transactions</a>
                {{ end }}
            {{ end }}
        </div>
    </div>
    {{ end }}

</main>

{{ template "footer.tpl" . }}