$(document).ready(function () {
    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);
        })
        .fail(function (exc) {
            updateDelegateUI(null);
        });

    $.when(getVotersRewards())
        .done(function (data) {
            updateVotersCountUI(data);
        })
        .fail(function (exc) {
            updateVotersCountUI(null);
        });

    $('#KryptoVotersRewards').DataTable({
        'processing': true,
        'ordering': true,
        'paging': true,
        'pageLength': 10,
        'bLengthChange': false,
        'order': [[3, 'desc']],
        'ajax': {
            'url': 'http://localhost:54000/voters/rewards',
            'dataSrc': 'data'
        },
        "dom": '<"top"f>rt<"bottom"l<"text-center"p>><"clear">',
        'columns': [
            { 'data': 'Address' },
            { 'data': 'Address' },
            { 'data': 'VoteWeight' },
            { 'data': 'EarnedAmountXX' },
            { 'data': 'VoteDuration' },
            { 'data': 'VoteDuration' }
        ],
        'columnDefs': [
            {
                'targets': [0],
                'searchable': false
            },
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
                    return '<span class="label label-primary label-block">' + getFormatedArkValue(data) + '</span>';
                },
                'targets': 3
            },
            {
                'render': function (data, type, row) {
                    return getFormatedArkDuration(data);
                },
                'targets': 4,
                'orderData': [5]
            },
            {
                'visible': false,
                'targets': 5
            }]
    });

    updateTableSearchBoxUI();
});