var barOptions = {
    series: {
        lines: {
            show: true,
            lineWidth: 2,
            fill: true
        },
        points: { show: true }
    },
    colors: ['#1ab394'],
    grid: {
        color: '#999999',
        hoverable: true,
        clickable: true,
        tickColor: '#D4D4D4',
        borderWidth: 0
    },
    legend: {
        show: false
    },
    alignTicksWithAxis: 1,
    tooltip: true,
    tooltipOpts: {
        content: 'x: %x, y: %y'
    },
    xaxis: {
        mode: "time",
        timeformat: "%d %b %Y"
    }
};

var chartSeries = [{
    data: [],
    lines: {
        fill: true
    }
}];

$(document).ready(function () {
    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);

            var urlParams = new URLSearchParams(window.location.search);
            if (urlParams.has('address')) {
                $('#tbAddress').val(urlParams.get('address'));
            }

            loadMyPaymentdUI();
        })
        .fail(function (exc) {
            updateDelegateUI(null);

            loadMyPaymentdUI();
        });

    $('#tbAddress').on('keyup', function (ev) {
        var keycode = ev.keyCode || ev.which;
        if (keycode == '13') {
            loadMyPaymentdUI();
        }
        return false;
    });

    $('#btnAddress').on('click', function() {
        loadMyPaymentdUI();
    });
});

var chart = $.plot($('#chartMyPayments'), chartSeries, barOptions);

var dt = $('#tblMyPayments').DataTable({
    'processing': true,
    'ordering': true,
    'searching': false,
    'paging': true,
    'pageLength': 10,
    'bLengthChange': false,
    'order': [[2, 'desc']],
    'columns': [
        { 'data': 'Address' },
        { 'data': 'EarnedAmountXX' },
        { 'data': 'CreatedAt' },
        { 'data': 'CreatedAt' },
        { 'data': 'Transaction.id' },
        { 'data': 'VoteDuration' },
        { 'data': 'VoteDuration' }
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
                return '<span class="label label-primary label-block">' + getFormatedArkValue(data) + '</span>';
            },
            'width': '15%',
            'targets': 1
        },
        {
            'render': function (data, type, row) {
                return getFormatedDateTime(data);
            },
            'width': '15%',
            'targets': 2,
            'orderData': [3]
        },
        {
            'render': function (data, type, row) {
                return getFormatedArkTransactionAddress(data);
            },
            'targets': 4
        },
        {
            'render': function (data, type, row) {
                return getFormatedArkDuration(data);
            },
            'targets': 5,
            'orderData': [6]
        },
        {
            'visible': false,
            'searchable': false,
            'targets': [3, 5]
        }
    ]
});

function reloadMyPaymentsGraph(data) {
    chartSeries[0].data = data;
    chart.setData(chartSeries);
    chart.setupGrid();
    chart.draw();
};

function reloadMyPaymentsTable(data) {
    dt.clear();
    dt.rows.add(data);
    dt.draw();
};

function loadMyPaymentdUI() {
    var address = $('#tbAddress').val();
    if (address.length == 0)
        return;

    $.when(getPaymentRunDetails(address))
        .done(function (pDetails) {
            var gData = [];
            if (pDetails.success) {
                var ordDetails = pDetails.data.sort(function (a, b) {
                    return moment(a.CreatedAt) > moment(b.CreatedAt) ? 1 : -1;
                });
                var gData = [];
                $.each(ordDetails, function (i, e) {
                    gData.push([moment(e.CreatedAt).valueOf(), e.EarnedAmountXX]);
                });

                var total = 0;
                $.each(ordDetails, function (i, e) {
                    total += e.EarnedAmountXX;
                });
                $('#lblTotalEarned').html('Ѧ ' + getFormatedArkValue(total));
                $('#lblNumOfTransactions').html(ordDetails.length);
            }
            else {
                $('#lblTotalEarned').html('Ѧ ERROR');
                $('#lblNumOfTransactions').html('ERROR');
            }
            reloadMyPaymentsGraph(gData);
            reloadMyPaymentsTable(pDetails.data);
        })
        .fail(function () {
            reloadMyPaymentsGraph([]);
            reloadMyPaymentsTable([]);
            $('#lblTotalEarned').html('Ѧ ERROR');
            $('#lblNumOfTransactions').html('ERROR');
        });
};