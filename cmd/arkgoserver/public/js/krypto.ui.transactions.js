$(document).ready(function () {
    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);
        })
        .fail(function (exc) {
            updateDelegateUI(null);
        });

    var qParams = new URLSearchParams(window.location.search);
    var meUrl = 'http://localhost:54000/delegate/paymentruns/details?parentid=' + qParams.get('pk');

    $('aMe').attr('href', meUrl);
    $('#lblTimestamp').html('(' + getFormatedDateTime(qParams.get('ts')) + ')');

    $('#KryptoPaymentTransactions').DataTable({
        'responsive': true,
        'processing': true,
        'ordering': true,
        'paging': true,
        'pageLength': 10,
        'bLengthChange': false,
        'ajax': {
            'url': meUrl,
            'dataSrc': 'data'
        },
        "dom": '<"top"f>rt<"bottom"l<"text-center"p>><"clear">',
        'columns': [
            { 'data': 'Address' },
            { 'data': 'VoteWeight' },
            { 'data': 'EarnedAmountXX' },
            { 'data': 'VoteDuration' },
            { 'data': 'VoteDuration' },
            { 'data': 'Transaction.id' }
        ],
        'columnDefs': [
            {
                'targets': [1, 2, 3],
                'searchable': false
            },
            {
                'render': function (data, type, row) {
                    return getFormatedArkAddress(data);
                },
                'targets': 0
            },
            {
                'render': function (data, type, row) {
                    return '<span class="label label-primary label-block">' + getFormatedArkValue(data) + '</span>';
                },
                'targets': 2
            },
            {
                'render': function (data, type, row) {
                    return getFormatedArkDuration(data);
                },
                'targets': 3,
                'orderData': [4]
            },
            {
                'visible': false,
                'targets': 4
            },
            {
                'render': function (data, type, row) {
                    return getFormatedArkTransactionAddress(data);
                },
                'targets': 5
            }]
    });

    updateTableSearchBoxUI();
});