{{ template "header.tpl" . }}

<script>
$(document).ready(function () {
    {{ range .assets }}
        addImage("{{ . }}");
    {{ end }}

    // Event delegation for all dynamically added images
    $(document).on('click', 'img[id^="dynamicImg_"]', function () {
        const clickedId = $(this).attr('id');
        $('#body').val(function(i, val) {
            return val + '\n![](' + clickedId.replace('dynamicImg_', '') + ')\n';
        });
        $('#body').focus();
    });

    function addImage(name) {
        var src = $(location).attr('origin') + '/web/assets/' + name;
        const imgId = 'dynamicImg_' + name;
        $('#assets').append('<div class="card" style="width: 18rem;"><img src="' + src + '" id="' + imgId + '" class="card-img-top"></div>');
    }

    $('#uploadBtn').on('click', function () {
        var fileInput = $('#imageUpload')[0];
        if (fileInput.files.length === 0) {
            alert('Please select an image.');
            return;
        }

        console.log("Uploading file: ", fileInput.files[0]);
        var formData = new FormData();
        formData.append('asset', fileInput.files[0]);

        $.ajax({
            url: 'upload',
            type: 'POST',
            data: formData,
            contentType: false,
            processData: false,
            success: function (response) {
                console.log("Upload successful, response: ", response);
                addImage(response);
                // Clear the file input
                $('#imageUpload').val('');
            },
            error: function () {
                alert('Upload failed.');
            }
        });
    });
});
</script>

<main>
    {{ with .Error }}
    <div class="alert alert-danger" role="alert">
        {{ . }}
    </div>
    {{ end }}

    <div class="row">
        <div class="col">
            <form action="/web/edit" method="POST">
                <input type="hidden" name="date" value="{{ .item.Date }}"/>
                <input type="hidden" id="user_id" value="{{ .UserID }}"/>

                <h5>{{ .item.Date }}</h5>

                <div class="mb-3">
                    <label for="title" class="form-label">Title:</label>
                    <input type="text" class="form-control" name="title" value="{{ .item.Title }}"/>
                </div>

                <div class="mb-3">
                    <label for="body" class="form-label">Body:</label>
                    <textarea class="form-control" name="body" id="body" rows="10">{{ .item.Body }}</textarea>
                </div>

                <div class="mb-3">
                    <label for="tags" class="form-label">Tags</label>
                    <input type="text" class="form-control" name="tags" value="{{ range .item.Tags }}{{.}}, {{ end }}"/>
                </div>

                <button type="submit" class="btn btn-primary">Save</button>
            </form>
        </div>
        <div class="col-3 text-center">`
            <input type="file" id="imageUpload" />
            <button id="uploadBtn" class="btn btn-secondary">Upload asset</button>

            <h5>Assets</h5>
            <div id="assets">
            </div>
        </div>
    </div>
</main>

{{ template "footer.tpl" . }}