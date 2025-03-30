{{ template "header.tpl" . }}

<main>
    <h2>Transactions</h2>
    <a href="/web/transactions?from={{.Last}}&accountID={{.AccountID}}" class="btn btn-primary" tabindex="-1" role="button">
        <i class="bi-arrow-left-circle-fill"></i>
    </a>
    {{ formatTime .From "2006-01-02" }} - {{ formatTime .To "2006-01-02" }}
    <a href="/web/transactions?from={{.Next}}&accountID={{.AccountID}}" class="btn btn-primary" tabindex="-1" role="button">
        <i class="bi-arrow-right-circle-fill"></i>
    </a>

    <form method="GET">
        <input type="hidden" name="from" value="{{.Current}}">
        <select name="accountID">
            <option value="">Select account</option>
            {{ range $.Accounts }}
            <option value="{{ .Id }}" {{ if eq .Id $.AccountID }}selected{{ end }}>
                {{.Type}}: {{ .Name }}
            </option>
            {{ end }}
        </select>
        <button class="btn btn-info" type="submit">Filter</button>
    </form>

    {{ range .Transactions }}
        <h5>{{ formatTime .Date "2006-01-02" }} {{ .Description }}</h5>
        <h6>Tags: 
            {{ range .Tags }}
                {{ . }}
            {{ end }}
        </h6>
        <p>
            {{ range .Movements }}
                {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
            {{ end }}
            <a href="/web/transactions/edit?id={{.ID}}" class="btn btn-primary" tabindex="-1" role="button">
                Edit
            </a>
        </p>
    {{ end }}

</main>

{{ template "footer.tpl" . }}
