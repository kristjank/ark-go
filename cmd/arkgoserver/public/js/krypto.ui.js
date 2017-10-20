function updateDelegateUI(data) {
    if (data == null || data.success == false) {
        $('#KryptoAddress').prop('href', '#');
        $('#KryptoUsername').html('ERROR');
        $('#lblUsername').html('ERROR');
        $('#KryptoBlocksForged').html('ERROR');
        $('#KryptoRate').html('ERROR');
    }
    else {
        $('#KryptoAddress').prop('href', getArkAddress(data.delegate.address));
        $('#KryptoUsername').html('delegate: ' + data.delegate.username);
        $('#lblUsername').html(data.delegate.username);
        $('#linkUsername').prop('href', getArkAddress(data.delegate.address));
        $('#KryptoBlocksForged').html(data.delegate.producedblocks);
        $('#KryptoRate').html(data.delegate.rate);
    }
};

function updateDelegateConfigUI(data) {
    if (data == null) {
        $('#lblShareRatio').html('ERROR');
        $('#lblFeesCovered').html('ERROR');
        $('#lblFeesCoveredInfo').html('ERROR');
        $('#lblMinAmount').html('ERROR');
        $('#lblFidelity').html('ERROR');
        $('#lblFidelityInfo').html('ERROR');
    }
    else {
        $('#lblShareRatio').html(data.shareratio * 100 + ' %');
        $('#lblFeesCovered').html(data.deductTxFees ? 'YES' : 'NO');
        $('#lblFeesCoveredInfo').html(data.deductTxFees ? 'Fees covered by voters' : 'Fees covered by delegate');
        $('#lblMinAmount').html(data.minamount + ' Ѧ');
        $('#lblFidelity').html(data.fidelityLimit);
        $('#lblFidelityInfo').html(data.fidelity ? 'ACTIVE' : 'INACTIVE');
        if (data.fidelityLimit <= 51) {
            $('#lblNodeForgingAllowed').html('FORGING');
            $('#widgetStatus').removeClass('red-bg').addClass('navy-bg');
        }
    }
};

function updateVotersCountUI(data) {
    if (data == null) {
        $('#KryptoVoters').html('ERROR');
    }
    else {
        $('#KryptoVoters').html(data.count);
    }
};

function updateTableSearchBoxUI() {
    $('.dataTables_filter').addClass('pull-left');
    $('.dataTables_paginate').css('text-align', 'center');
};

function updateSocialFeedUI(data, pagingSize) {
    pagingSize = pagingSize || 5;

    data.transactions.sort(function(a, b) { return b.timestamp - a.timestamp });

    $('#listSocialFeed').html(updateSocialFeedPageUI(data, 0, pagingSize));
    $('#paginationSocialFeed').html(updateFeedPaginationUI(data, 0, pagingSize));

    updateLinkifyUI();

    $('#paginationSocialFeed').on('click', '.paginate_button', function () {
        var me = $(this);
        if (!me.hasClass('disabled') && !me.hasClass('active')) {
            var p = parseInt(me.attr('data-page'));

            $('#listSocialFeed').html(updateSocialFeedPageUI(data, p, pagingSize));
            $('#paginationSocialFeed').html(updateFeedPaginationUI(data, p, pagingSize));

            updateLinkifyUI();
        }
        return false;
    });
};

function updateSocialFeedPageUI(data, page, pageSize) {
    var html = '';

    if (data != null && data.success == true) {
        var i = page * pageSize;
        var j = 0;
        while (j < pageSize && i < data.transactions.length) {
            var e = data.transactions[i];

            html += '\
                <li class="list-group-item">\
                    <p>' + e.vendorField + '</p>\
                    <small class="block text-muted"><i class="fa fa-clock-o"></i> ' + getTimeCounter(e.timestamp) + '</small>\
                </li>';

            i++;
            j++;
        }
    }

    return html;
};

function updateFeedPaginationUI(data, page, pageSize) {
    if (data == null || data.transactions.length == 0)
        return '';

    var maxPage = Math.ceil(data.transactions.length / pageSize);

    var html = '<ul class="pagination">';
    html += '<li class="paginate_button previous' + (page == 0 ? ' disabled' : '') + '"' + (page == 0 ? '' : ' data-page="' + (page - 1) + '"') + '><a href="#">Previous</a></li>';

    if (maxPage <= 7) {
        for (var i = 0; i < maxPage; i++) {
            html += '<li class="paginate_button' + (i == page ? ' active' : '') + '" data-page="' + i + '"><a href="#">' + (i + 1) + '</a></li>';
        }
    }
    else {
        if (page <= 4) {
            for (var i = 0; i < 5; i++) {
                html += '<li class="paginate_button' + (i == page ? ' active' : '') + '" data-page="' + i + '"><a href="#">' + (i + 1) + '</a></li>';
            }
            html += '<li class="paginate_button disabled"><a href="#">…</a></li>';
        }
        else if (page > (maxPage - 5)) {
            html += '<li class="paginate_button" data-page="0"><a href="#">1</a></li>';
            html += '<li class="paginate_button disabled"><a href="#">…</a></li>';
            for (var i = (maxPage - 5); i < maxPage; i++) {
                html += '<li class="paginate_button' + (i == page ? ' active' : '') + '" data-page="' + i + '"><a href="#">' + (i + 1) + '</a></li>';
            }
        }
        else {
            html += '<li class="paginate_button" data-page="0"><a href="#">1</a></li>';
            html += '<li class="paginate_button disabled"><a href="#">…</a></li>';

            for (var i = (page - 1); i <= (page + 1); i++) {
                html += '<li class="paginate_button' + (i == page ? ' active' : '') + '" data-page="' + i + '"><a href="#">' + (i + 1) + '</a></li>';
            }

            html += '<li class="paginate_button disabled"><a href="#">…</a></li>';
            html += '<li class="paginate_button" data-page="' + (maxPage - 1) + '"><a href="#">' + maxPage + '</a></li>';
        }
    }

    html += '<li class="paginate_button next' + (page == (maxPage - 1) ? ' disabled' : '') + '"' + (page == (maxPage - 1) ? '' : ' data-page="' + (page + 1) + '"') + '><a href="#">Next</a></li>';
    return html;
};

function updateSocialInfoUI(data) {
    if (data == null) {
        $('#colSocialInfo').html('');
    }
    else {
        data.twitter == undefined || data.twitter.length == 0
            ? $('#btnTwitter').remove()
            : $('#btnTwitter').prop('href', 'https://twitter.com/' + data.twitter);

        data.reddit == undefined || data.reddit.length == 0
            ? $('#btnReddit').remove()
            : $('#btnReddit').prop('href', 'https://www.reddit.com/user/' + data.reddit);

        data.slack == undefined || data.slack.length == 0
            ? $('#btnSlack').remove()
            : $('#btnSlack').prop('href', 'slack://user?team={TEAM_ID}&id=' + data.slack);

        data.proposal == undefined || data.proposal.length == 0
            ? $('#btnProposal').remove()
            : $('#btnProposal').prop('href', '' + data.proposal);

        data.github == undefined || data.github.length == 0
            ? $('#btnGithub').remove()
            : $('#btnGithub').prop('href', 'https://github.com/' + data.github);

        data.email == undefined || data.email.length == 0
            ? $('#btnEmail').remove()
            : $('#btnEmail').prop('href', 'mailto:' + data.email);
    }
};

function updateNodeStatusUI(data) {
    if (data == null || data.success == false) {
        $('#widgetStatus').addClass('red-bg');
        $('#lblNodeForgingAllowed').html('FORGING UNDEFINED');
        $('#lblNodeHeight').html('Height undefined');
    }
    else {
        $('#widgetStatus').removeClass('navy-bg red-bg');
        $('#widgetStatus').addClass(data.forgingAllowed ? 'navy-bg' : 'red-bg');
        $('#lblNodeForgingAllowed').html(data.forgingAllowed ? 'FORGING' : 'NOT FORGING');
        $('#lblNodeHeight').html('Height ' + data.height);
    }
};

function updateLinkifyUI() {
    $('#listSocialFeed').linkify({
        target: "_blank"
    });
};

function getTimeCounter(t) {
    var arkStart = moment([2017, 2, 21, 13, 0, 0, 0]);
    var dt = arkStart.add(t, 's');

    var now = moment();

    return dt.format('LLL') + ' (' + dt.to(now, true) + ' ago)';
}

function getArkAddress(address) {
    return 'https://explorer.ark.io/address/' + address;
};

function getFormatedArkValue(value) {
    return parseFloat(value).toFixed(8);
};

function getFormatedArkAddress(address) {
    return '<a href="https://explorer.ark.io/address/' + address + '" target="_blank">' + address + '</a>';
};

function getFormatedArkTransactionAddress(tid) {
    return '<a href="https://explorer.ark.io/tx/' + tid + '" target="_blank">View Tx</a>';
};

function getFormatedArkDuration(value) {
    var res = '<span class="label label-info">' + value + ' h</span>';

    var nDays = value / 24;
    var nYears = nDays / 365;

    if (nYears > 1) {
        res += ' <span class="label">over a year</span>';
    }
    else {
        var nMonths = nDays / 30;

        if (nMonths > 2)
            res += ' <span class="label">around ' + Math.floor(nMonths) + ' months</span>';
        else if (nMonths > 1)
            res += ' <span class="label">over a month</span>';
        else
            res += ' <span class="label">around ' + Math.floor(nDays) + ' days</span>';
    }

    return res;
};

function getFormatedDateTime(value) {
    return moment(value).format('YYYY-MMM-DD HH:mm');
};

$(document).ready(function () {
    updateDatatblesResizingBugUI();
    $(window).resize(function() {
        updateDatatblesResizingBugUI();
    });
});

function updateDatatblesResizingBugUI() {
    $('.dt-responsive').css('width', '100%');
};