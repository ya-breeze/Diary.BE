{{ template "header.tpl" . }}

<main>
    <h2>Create matcher</h2>
    <form method="POST">
        <input type="hidden" name="id" value="{{ .Matcher.Id }}">
        <input type="hidden" name="transaction_id" value="{{ .Transaction.ID }}">

        <div class="mb-3">
            <label for="name" class="form-label">Name</label>
            <input type="text" class="form-control" name="name" value="{{ .Matcher.Name }}">
        </div>

        <div class="mb-3">
            <label for="outputDescription" class="form-label">Output description</label>
            <input type="text" class="form-control" name="outputDescription" value="{{ .Matcher.OutputDescription }}">
        </div>
        <div class="mb-3">
            <label for="outputTags" class="form-label">Output tags</label>
            <input type="text" class="form-control" name="outputTags" value="{{ range .Matcher.OutputTags }}{{.}}, {{ end }}">
        </div>

        <div class="mb-3">
            <label for="descriptionRegExp" class="form-label">Description regexp</label>
            <input type="text" class="form-control" name="descriptionRegExp" value="{{ .Matcher.DescriptionRegExp }}">
        </div>
        <div class="mb-3">
            <label for="partnerAccountNumberRegExp" class="form-label">Partner account regexp</label>
            <input type="text" class="form-control" name="partnerAccountNumberRegExp" value="{{ .Matcher.PartnerAccountNumberRegExp }}">
        </div>
        <div class="mb-3">
            <label for="account" class="form-label">Account</label>
            <select class="form-select" name="account">
                <option value="">Select account</option>
                {{ range $.Accounts }}
                <option value="{{ .Id }}" {{ if eq .Id $.Matcher.OutputAccountId }}selected{{ end }}>
                    {{.Type}}: {{ .Name }}
                </option>
                {{ end }}
            </select>
        </div>

        <button class="btn btn-primary" type="submit">Save</button>
    </form>
</main>

{{ template "footer.tpl" . }}