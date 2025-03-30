{{ template "header.tpl" . }}

<main>
    {{ with .Unprocessed }}
        <h5></h5>
        <form class="row g-3" action="/web/unprocessed/convert" method="POST">
            <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">

            <strong><i class="bi-calendar-event"></i> Date: {{ formatTime .Transaction.Date "2006-01-02" }}</strong>

            <div class="mb-3">
                <label for="description" class="form-label">Description: </label>
                <input type="text" class="form-control" name="description" value="{{ .Transaction.Description }}">
            </div>

            <div class="mb-3">
                <i class="bi-bank"></i> Movements:
                {{ range $i, $m := .Transaction.Movements }}
                    <div class="mb-3">
                        <label for="account" class="form-label">{{ $m.Amount }} {{ $m.CurrencyName }}</label>
                        <select class="form-select" name="account_{{ $i }}">
                            <option value="">Select account</option>
                            {{ range $.Accounts }}
                            <option value="{{ .Id }}" {{ if eq .Id $m.AccountID }}selected{{ end }}>
                                {{.Type}}: {{ .Name }}
                            </option>
                            {{ end }}
                        </select>
                    </div>
                {{ end }}
            </p>

            <div class="mb-3">
                {{ if ne .Transaction.PartnerName "" }}
                    Partner name: {{ .Transaction.PartnerName }}
                {{ end }}
            </div>
            <div class="mb-3">
                {{ if ne .Transaction.PartnerAccount "" }}
                    <i class="bi-bank"></i> Partner account: {{ .Transaction.PartnerAccount }}
                {{ end }}
            </div>
            <div class="mb-3">
                {{ if ne .Transaction.Place "" }}
                    <i class="bi-pin-map-fill"></i> Place: {{ .Transaction.Place }}
                {{ end }}
            </div>
            <div class="mb-3">
                <i class="bi-tags"></i> Tags: 
                {{ range .Transaction.Tags }}
                    {{ . }}
                {{ end }}
            </div>
            

            <div class="mb-3">
                <button class="btn btn-success btn-lg" type="submit">Convert</button>
                <a href="/web/matchers/edit?transaction_id={{ .Transaction.ID }}" class="btn btn-outline-primary pull-right">Create matcher</a>
                <a href="?id={{ .Transaction.ID }}" class="btn btn-outline-secondary pull-right">Skip ({{ decrease $.UnprocessedCount }} left)</a>
            </div>
        </form>

        {{ with .Matched }}
            <h3>Matched</h3>
            {{ range . }}
                <form action="/web/unprocessed/convert" method="POST">
                    <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">
                    {{ range $i, $m := .Transaction.Movements }}
                        <input type="hidden" name="account_{{ $i }}" value="{{ $m.AccountID }}">
                    {{ end }}
                    <button type="submit">Convert</button>
                </form>
                <h5>{{ formatTime .Transaction.Date "2006-01-02" }} {{ .Transaction.Description }}</h5>
                {{ if ne .Transaction.PartnerName "" }}
                    <h6>Partner name: {{ .Transaction.PartnerName }}</h6>
                {{ end }}
                {{ if ne .Transaction.PartnerAccount "" }}
                    <h6>Partner account: {{ .Transaction.PartnerAccount }}</h6>
                {{ end }}
                <h6>Tags: 
                {{ range .Transaction.Tags }}
                    {{ . }},
                {{ end }}
                </h6>
                <p>
                    {{ range .Transaction.Movements }}
                        {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
                    {{ end }}
                </p>
            {{ end }}
        {{ end }}


        {{ if .Duplicates }}
            <h3>Duplicates</h3>
            {{ range .Duplicates }}
                <h5>{{ formatTime .Date "2006-01-02" }} {{ .Description }}</h5>
                <p>
                    {{ range .Movements }}
                        {{ .Amount }} {{ .CurrencyName }} [{{ .AccountName }}] <br>
                    {{ end }}
                </p>
                <a href="/web/unprocessed/delete?id={{$.Unprocessed.Transaction.ID}}&duplicateOf={{.ID}}" class="btn btn-primary" tabindex="-1" role="button">
                    Mark as duplicate
                </a>
            {{ end }}
        {{ end }}

    {{ end }}

</main>

{{ template "footer.tpl" . }}