function showTime() {
    var now = moment();

    $('[data-module=clock]').each(function() {
        var $this = $(this),
            settings = $this.data('settings');

        var t = now;

        if (settings && settings.timezone) {
            t = t.clone().tz(settings.timezone);
        }

        $this.find('.clock-large').text(t.format('LTS'));
        $this.find('.date-large').text(t.format('dddd, LL'));
    });
}

setInterval(showTime, 500);

showTime();

window.modules.clock = {
    settings: function(values) {
        var options = [];

        for (var key in moment.tz._names) {
            var value = moment.tz._names[key];

            options.push({
                text: value,
                value: value
            });
        }

        return {
            timezone : {
                label: 'Timezone',
                type: 'select',
                options: options
            }
        }
    }
};