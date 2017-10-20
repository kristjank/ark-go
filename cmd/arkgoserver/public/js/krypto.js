var API_URL_BASE = 'http://localhost:54000/';
var API_URL_DELEGATE = API_URL_BASE + 'delegate';
var API_URL_DELEGATE_CONFIG = API_URL_BASE + 'delegate/config';
var API_URL_PAYMENTRUNS = API_URL_BASE + 'delegate/paymentruns';
var API_URL_PAYMENTRUNS_DETAILS = API_URL_PAYMENTRUNS + '/details';
var API_URL_VOTERS = API_URL_BASE + 'voters';
var API_URL_VOTERS_BLOCKED = API_URL_BASE + 'voters/blocked';
var API_URL_VOTERS_REWARDS = API_URL_BASE + 'voters/rewards';
var API_URL_SOCIAL = API_URL_BASE + 'social';
var API_URL_SOCIAL_INFO = API_URL_SOCIAL + '/info';
var API_URL_NODESTATUS = API_URL_BASE + 'delegate/nodestatus';

function baseGet(url, onSuccess, onFail) {
    $.ajax({
        url: url,
        type: 'GET',
        cache: false,
        contentType: 'application/json; charset-utf-8',
        beforeSend: function (xhr, settings) {
            console.debug('AJAX GET: "' + url + '"');
        },
        success: function (data) {
            if (onSuccess != null)
                onSuccess(data);
        },
        error: function (jqXHR, textStatus, errorThrown) {
            if (onFail != null)
                onFail(jqXHR, textStatus, errorThrown);
        }
    });
};

function getDelegate() {
    var t = $.Deferred();

    baseGet(
        API_URL_DELEGATE,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire DELEGATE'
            });
        });

    return t.promise();
};

function getDelegateConfig() {
    var t = $.Deferred();

    baseGet(
        API_URL_DELEGATE_CONFIG,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire DELEGATE CONFIG'
            });
        });

    return t.promise();
};

function getPaymentRuns() {
    var t = $.Deferred();

    baseGet(
        API_URL_PAYMENTRUNS,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire PAYMENT RUNS'
            });
        });

    return t.promise();
};

function getPaymentRunDetails(address) {
    var t = $.Deferred();

    baseGet(
        API_URL_PAYMENTRUNS_DETAILS + '?address=' + address,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire PAYMENT RUN DETAILS'
            });
        });

    return t.promise();
};

function getVotersBlocked() {
    var t = $.Deferred();
    
        baseGet(
            API_URL_VOTERS_BLOCKED,
            function (data) { t.resolve(data); },
            function (jqXHR) {
                t.reject({
                    Exception: jqXHR,
                    Message: 'Failed to acquire VOTERS BLOCKED'
                });
            });
    
        return t.promise();
};

function getVotersActive() {
    var t = $.Deferred();
    
        baseGet(
            API_URL_VOTERS,
            function (data) { t.resolve(data); },
            function (jqXHR) {
                t.reject({
                    Exception: jqXHR,
                    Message: 'Failed to acquire VOTERS ACTIVE'
                });
            });
    
        return t.promise();
};

function getVotersRewards() {
    var t = $.Deferred();

    baseGet(
        API_URL_VOTERS_REWARDS,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire VOTERS REWARDS'
            });
        });

    return t.promise();
};

function getSocialFeed() {
    var t = $.Deferred();

    baseGet(
        API_URL_SOCIAL,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire SOCIAL'
            });
        });

    return t.promise();
};

function getSocialInfo() {
    var t = $.Deferred();

    baseGet(
        API_URL_SOCIAL_INFO,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire SOCIAL INFO'
            });
        });

    return t.promise();
};

function getNodeStatus() {
    var t = $.Deferred();

    baseGet(
        API_URL_NODESTATUS,
        function (data) { t.resolve(data); },
        function (jqXHR) {
            t.reject({
                Exception: jqXHR,
                Message: 'Failed to acquire NODE STATUS'
            });
        });

    return t.promise();
};

function getRedditFeed(channel) {
    var t = $.Deferred();

    var url = 'https://www.reddit.com/r/' + channel + '/.json?callback=';
        $.ajax({
            url: url,
            type: 'GET',
            cache: false,
            contentType: 'application/json; charset-utf-8',
            beforeSend: function (xhr, settings) {
                console.debug('AJAX GET: "' + url + '"');
            },
            success: function (data) {
                t.resolve(data);
            },
            error: function (jqXHR, textStatus, errorThrown) {
                t.reject({
                    Exception: jqXHR,
                    Message: 'Failed to acquire REDDIT FEED'
                });
            }
        });

    return t.promise();
};