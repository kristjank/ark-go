$(document).ready(function () {
    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);
        })
        .fail(function (exc) {
            updateDelegateUI(null);
        });

    $('#tblKryptoPayments').DataTable({
        'processing': true,
        'ordering': true,
        'searching': false,
        'paging': true,
        'pageLength': 10,
        'bLengthChange': false,
        'order': [[0, 'desc']],
        'ajax': {
            'url': 'http://localhost:54000/delegate/paymentruns',
            'dataSrc': 'data'
        },
        'columns': [
            { 'data': 'CreatedAt' },
            { 'data': 'CreatedAt' },
            { 'data': 'ShareRatio' },
            { 'data': 'Fidelity' },
            { 'data': 'FeeDeduction' },
            { 'data': 'NrOfTransactions' },
            { 'data': 'Pk' }
        ],
        'columnDefs': [
            {
                'render': function (data, type, row) {
                    return getFormatedDateTime(data);
                },
                'targets': 0,
                'orderData': [1]
            },
            {
                'visible': false,
                'targets': 1
            },
            {
                'render': function (data, type, row) {
                    return '<a class="btn btn-block btn-xs btn-primary" href="transactions.html?pk=' + data + '&ts=' + row.CreatedAt + '" target="blank">DETAILS</a>';
                },
                'targets': 6
            }
        ]
    });
});