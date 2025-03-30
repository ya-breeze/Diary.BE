{{ template "header.tpl" . }}

<main>
    <h2>Accounts</h2>
    {{ with .Error }}
    <div class="alert alert-danger" role="alert">
        {{ . }}
    </div>
    {{ end }}

    <div class="row">
        <div class="col-6">
            <form method="POST">
                <input type="hidden" name="id" value="{{ .Id }}">
                <div class="form-group">
                    <label for="name">Name</label>
                    <input type="text" class="form-control" id="name" name="name" value="{{.Name}}">
                </div>
                <div class="form-group">
                    <label for="type">Type</label>
                    <select class="form-select" aria-label="Account type" name="type">
                        <option value="expense" {{ if eq .Type "expense"}}selected{{ end }}>Expense</option>
                        <option value="asset" {{ if eq .Type "asset"}}selected{{ end }}>Asset</option>
                        <option value="income" {{ if eq .Type "income"}}selected{{ end }}>Income</option>
                    </select>
                </div>
                <div class="form-group">
                    <label for="description">Description</label>
                    <input type="text" class="form-control" id="description" name="description"
                        value="{{.Description}}">
                </div>
                <button type="submit" class="btn btn-primary">Save account</button>
            </form>
        </div>
    </div>
</main>

{{ template "footer.tpl" . }}