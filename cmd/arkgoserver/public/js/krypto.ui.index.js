var chartOptions = {
    colors: ['#1ab394', '#941ab3'],
    grid: {
        color: '#999999',
        hoverable: false,
        clickable: false,
        tickColor: '#D4D4D4',
        borderWidth: 0
    },
    legend: {
        labelBoxBorderColor: '#EEEEEE',
        position: 'nw'
    },
    alignTicksWithAxis: 1,
    tooltip: true,
    tooltipOpts: {
        content: 'x: %x, y: %y'
    },
    xaxis: {
        mode: 'time',
        timeformat: '%d %b %Y'
    },
    yaxes: [
        {
            show: true,
            axisLabel: 'Weight',
            points: { show: false },
            axisLabelUseCanvas: true,
            tickFormatter: function (v, axis) {
                return v + ' Ѧ';
            }
        },
        {
            show: true,
            axisLabel: 'Number',
            color: '#941ab3',
            position: 'right',
            axisLabelUseCanvas: true
        }
    ]
};

var chartSeries = [{
    label: 'Vote Height (Ѧ)',
    lines: { show: true, fill: true, fillColor: "rgba(26, 179, 148, 0.2)" },
    data: []
}, {
    label: 'Num. of Transactions',
    lines: { show: true, fill: true, fillColor: "rgba(148, 26, 179, 0.1)" },
    data: [],
    yaxis: 2,
}];

$(document).ready(function () {
    $.when(getNodeStatus())
        .done(function (data) {
            updateNodeStatusUI(data);
        })
        .fail(function (exc) {
            updateNodeStatusUI(null);
        });

    $.when(getDelegate())
        .done(function (data) {
            updateDelegateUI(data);
        })
        .fail(function (exc) {
            updateDelegateUI(null);
        });

    loadMyPaymentdUI();

    $.when(getVotersRewards())
        .done(function (data) {
            updateVotersCountUI(data);
        })
        .fail(function (exc) {
            updateVotersCountUI(null);
        });

    $.when(getSocialFeed())
        .done(function (data) {
            updateSocialFeedUI(data);
        })
        .fail(function (exc) {
            updateSocialFeedUI(null);
        });
});

var chart = $.plot($('#chartMyPayments'), chartSeries, chartOptions);

function reloadMyPaymentsGraph(data1, data2) {
    chartSeries[0].data = data1;
    chartSeries[1].data = data2;
    chart.setData(chartSeries);
    chart.setupGrid();
    chart.draw();
};

function loadMyPaymentdUI() {
    $.when(getPaymentRuns())
        .done(function (pDetails) {
            var gData = [];
            if (pDetails.success) {
                var ordDetails = pDetails.data.sort(function (a, b) {
                    return moment(a.CreatedAt) > moment(b.CreatedAt) ? 1 : -1;
                });
                var gData = [];
                var hData = [];
                $.each(ordDetails, function (i, e) {
                    gData.push([moment(e.CreatedAt).valueOf(), getFormatedArkValue(e.VoteWeight / Math.pow(10, 8))]);
                    hData.push([moment(e.CreatedAt).valueOf(), e.NrOfTransactions]);
                });
            }
            reloadMyPaymentsGraph(gData, hData);
        })
        .fail(function () {
            reloadMyPaymentsGraph([], []);
        });
};