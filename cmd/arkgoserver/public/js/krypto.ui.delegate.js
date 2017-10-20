$(document).ready(function () {
    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);
        })
        .fail(function (exc) {
            updateDelegateUI(null);
        });

    $.when(getDelegateConfig())
        .done(function (data) {
            updateDelegateConfigUI(data);
        })
        .fail(function (exc) {
            updateDelegateConfigUI(null);
        });

    $.when(getSocialInfo())
        .done(function (data) {
            updateSocialInfoUI(data);
        })
        .fail(function (exc) {
            updateSocialInfoUI(null);
        });

    $('#tblBlockedVoters').DataTable({
        'processing': true,
        'ordering': true,
        'searching': false,
        'paging': true,
        'pageLength': 10,
        'bLengthChange': false,
        'order': [[1, 'desc']],
        'ajax': {
            'url': 'http://localhost:54000/voters/blocked'
        },
        'columns': [
            { 'data': 'address' },
            { 'data': 'balance' }
        ],
        'columnDefs': [
            {
                'render': function (data, type, row) {
                    return getFormatedArkAddress(data);
                },
                'targets': 0
            },
            {
                'render': function (data, type, row) {
                    return '<span class="label label-warning">' + data + '</span>';
                },
                'targets': 1
            }
        ]
    });

    $('#tblVoters').DataTable({
        'processing': true,
        'ordering': true,
        'searching': false,
        'paging': true,
        'pageLength': 10,
        'bLengthChange': false,
        'order': [[2, 'desc']],
        'ajax': {
            'url': 'http://localhost:54000/voters',
            'dataSrc': 'accounts'
        },
        'columns': [
            { 'data': 'address' },
            { 'data': 'address' },
            { 'data': 'balance' }
        ],
        'columnDefs': [
            {
                'render': function (data, type, row) {
                    return '<a href="yourpayments.html?address=' + data + '" target="_blank">See Payments</a>';
                },
                'targets': 0,
                'sortable': false,
                'width': '10%'
            },
            {
                'render': function (data, type, row) {
                    return getFormatedArkAddress(data);
                },
                'targets': 1
            },
            {
                'render': function (data, type, row) {
                    return '<span class="label label-primary label-block">' + getFormatedArkValue(parseInt(data) / Math.pow(10, 8)) + '</span>';
                },
                'targets': 2
            }
        ]
    });
});
