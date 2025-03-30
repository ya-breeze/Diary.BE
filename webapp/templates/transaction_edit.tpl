{{ template "header.tpl" . }}

<main>
    <h2>Transaction</h2>
    {{ with .Error }}
    <div class="alert alert-danger" role="alert">
        {{ . }}
    </div>
    {{ end }}

    <form action="/web/transactions/edit" method="POST">
        <input type="hidden" name="id" value="{{ .Transaction.ID }}">
        <h5>{{ formatTime .Transaction.Date "2006-01-02" }}</h5>
        <div class="mb-3">
            <label for="description" class="form-label">Description:</label>
            <input type="text" class="form-control" name="description" value="{{ .Transaction.Description }}">
        </div>

        {{ if ne .Transaction.PartnerName "" }}
        <h6>Partner name: {{ .Transaction.PartnerName }}</h6>
        {{ end }}
        {{ if ne .Transaction.PartnerAccount "" }}
        <h6>Partner account: {{ .Transaction.PartnerAccount }}</h6>
        {{ end }}
        {{ if ne .Transaction.Place "" }}
        <h6>Place: {{ .Transaction.Place }}</h6>
        {{ end }}
        <div class="mb-3">
            <label for="tags" class="form-label">Tags</label>
            <input type="text" class="form-control" name="tags" value="{{ range .Transaction.Tags }}{{.}}, {{ end }}">
        </div>
        </h6>
        <p>
            <h6>Movements</h6>
            {{ range $i, $m := .Transaction.Movements }}
                <div class="card" style="width: 18rem;">
                    <input type="hidden" name="currency_{{ $i }}" value="{{ $m.CurrencyID }}">
                    <div class="card-body">
                        <div class="mb-3">
                            <label for="amount_{{ $i }}" class="form-label">Amount ({{ $m.CurrencyName }})</label>
                            <input type="text" class="form-control" name="amount_{{ $i }}" value="{{ $m.Amount }}">
                        </div>

                        <div class="mb-3">
                            <label for="account" class="form-label">Account</label>
                            <select class="form-select" name="account_{{ $i }}">
                                <option value="">Select account</option>
                                {{ range $.Accounts }}
                                <option value="{{ .Id }}" {{ if eq .Id $m.AccountID }}selected{{ end }}>
                                    {{.Type}}: {{ .Name }}
                                </option>
                                {{ end }}
                            </select>
                        </div>
                    </div>
                </div>
            {{ end }}
        </p>
        <button type="submit">Save</button>
    </form>
</main>

{{ template "footer.tpl" . }}