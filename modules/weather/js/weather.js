window.modules.weather = {
    settings : function(values) {
        return {
            location: {
                label: 'Location',
                type: 'text',
                value: values.location || ''
            },
            refresh: {
                label: 'Refresh Time',
                type: 'select',
                value: values.refresh || '60',
                options: [
                    { text : '1 minute', value : '60' },
                    { text : '5 minutes', value : '300'},
                    { text : '10 minutes', value : '600' }
                ]
            }
        };
    },

    saveSettings : function(values, success) {
        $.request('weather::onGeocoding', {
            data: { location : values.location },
            success: function(data) {
                var location = data.results[0].geometry.location;

                success({
                    location: data.results[0].formatted_address,
                    latitude: location.lat,
                    longitude: location.lng,
                    refresh: values.refresh
                });
            }
        });
    }
};