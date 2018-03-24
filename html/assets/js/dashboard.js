$(function () {
    window.modules = {};

    $('.grid-stack').gridstack({
        float: true
    });

    var counter = 0;

    $.request('onLoadLayout', {
        success: function(data) {
            for (var i = 0; i < data.length; i++) {
                addWidget(data[i]);
            }

            initializeModules();
        }
    });

    function addWidget(widget) {
        var grid = $('.grid-stack').data('gridstack');

        var $elem = $('<div><div class="grid-stack-item-content" /></div>');

        var $editElem = $('<a href="#" class="grid-edit"><i class="fa fa-cog"></i></a>'),
            $removeElem = $('<a href="#" class="grid-remove"><i class="fa fa-times"></i></a>');

        $elem.append($('<div/>').addClass('grid-module-edit').append($editElem, '&nbsp;', $removeElem));

        $elem.addClass('module-' + widget.module)
            .attr('id', widget.module + '-' + counter)
            .attr('data-module', widget.module)
            .attr('data-settings', widget.settings)
            .attr('data-request', widget.module + '::onRender');

        counter++;

        grid.addWidget($elem, widget.x, widget.y, widget.width, widget.height);

        return $elem;
    }

    function initializeModules() {
        $('[data-module]').each(function() {
            var $this = $(this);

            initializeModule($this);
        });
    }

    function initializeModule($this) {
        $this.on('moduleRefresh', function() {
            $this.request($this.data('request'), {
                data : $.extend({ id : $this.attr('id') }, $this.data('settings'))
            });
        });

        $this.trigger('moduleRefresh');

        var settings = $this.data('settings');

        if (settings && settings.refresh) {
            $this.data('interval', setInterval(function() {
                $this.trigger('moduleRefresh');
            }, settings.refresh * 1000));
        }
    }

    $(document).on('click', '.grid-edit', function() {
        var $this = $(this),
            $parent = $this.closest('.grid-stack-item'),
            module = $parent.data('module');

        var fields = {};

        if (module in window.modules && 'settings' in window.modules[module]) {
            fields = window.modules[module].settings($parent.data('settings'));
        }

        bootbox.form({
            title: 'Module Settings',
            fields: fields,
            callback: function (values) {
                if (!values) {
                    return;
                }

                // Go to weather's script to validate data, etc
                if (module in window.modules && 'saveSettings' in window.modules[module]) {
                    window.modules[module].saveSettings(values, function(finalValues) {
                        $parent.data('settings', finalValues);
                        $parent.trigger('moduleRefresh');

                        saveWidgetLayout();
                    });
                } else {
                    $parent.data('settings', values);
                    $parent.trigger('moduleRefresh');

                    saveWidgetLayout();
                }
            }
        });
    });

    $(document).on('click', '.grid-remove', function() {
        var gridstack = $('.grid-stack').data('gridstack'),
            $item = $(this).closest('.grid-stack-item');

        gridstack.removeWidget($item.get(0));

        var interval = $item.data('interval');

        if (interval) {
            window.clearInterval(interval);
        }

        saveWidgetLayout();
    });

    Mousetrap.bind(['command+s', 'ctrl+s'], function() {
        saveWidgetLayout();
        return false;
    });

    Mousetrap.bind(['command+a', 'ctrl+a'], function() {
        $.request('onListWidgets', {
            success: function(data) {
                var fields = {
                    widget : {
                        label: 'Widget',
                        type: 'select',
                        options: []
                    }
                };

                for (var id in data) {
                    fields.widget.options.push({ text : data[id].name, value : id });
                }

                bootbox.form({
                    title: 'Add Widget',
                    fields: fields,
                    callback: function (values) {
                        var $elem = addWidget({
                            module : values.widget,
                            settings: {},
                            width: data[values.widget].defaultSize[0],
                            height: data[values.widget].defaultSize[1],
                            x: 0,
                            y: 0
                        });

                        $elem.find('.grid-edit').click();

                        initializeModule($elem);
                    }
                });
            }
        });
        return false;
    });

    function saveWidgetLayout() {
        var serializedData = _.map($('.grid-stack > .grid-stack-item:visible'), function (el) {
            el = $(el);
            var node = el.data('_gridstack_node');

            return {
                module: el.data('module'),
                settings: JSON.stringify(el.data('settings')),
                x: node.x,
                y: node.y,
                width: node.width,
                height: node.height
            };
        }, this);

        $.request('onSaveLayout', {
            data: {
                layout : JSON.stringify(serializedData)
            },
            success : function(data) {
                this.success(data);

                $.notify('Successfully saved layout.', 'success');
            }
        });
    }
});