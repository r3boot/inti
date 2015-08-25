/*jslint browser: true*/
/*global $, jQuery, alert*/
"use strict";

var dmx_config = [];
var color_wheel = null;

function get_json(url) {
    return JSON.parse($.ajax({
        url: url,
        type: 'GET',
        dataType: 'json',
        global: false,
        async: false,
        success: function(data) {
            return data;
        }
    }).responseText);
}


function get_fixture_byaddr(bus_id, address) {
    var bus_name, bus, f_name, fixture = null;

    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            if (bus.id === bus_id) {
                for (f_name in bus.fixtures) {
                    if (bus.fixtures.hasOwnProperty(f_name)) {
                        fixture = bus.fixtures[f_name];
                        if (fixture.address === address) {
                            return fixture;
                        }
                    }
                }
            }
        }
    }

    return;
}


function transfer_buffer(data) {
    $.ajax({
        url: '/v1/buffer',
        type: 'POST',
        data: JSON.stringify(data),
        dataType: 'json',
    });
}


function set_fixture_color(color) {
    var bus_name, bus, f_selector, tokens, f_addr = null;
    var fixture, red_offset, green_offset, blue_offset, address = null;

    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            f_selector = 'f_bus-' + bus.id;
            $('#'+f_selector+' :selected').each(function(i, selected) {
                tokens = selected.value.split('-');
                f_addr = parseInt(tokens[tokens.length - 1]);
                fixture = get_fixture_byaddr(bus.id, f_addr);
                red_offset = fixture.channels.red;
                green_offset = fixture.channels.green;
                blue_offset = fixture.channels.blue;

                bus.buffer[f_addr + red_offset] = parseInt(color.r);
                bus.buffer[f_addr + green_offset] = parseInt(color.g);
                bus.buffer[f_addr + blue_offset] = parseInt(color.b);
            });

            transfer_buffer(bus);
        }
    }
}


function set_fixture_pan(value) {
    var bus_name, bus, f_selector, tokens, f_addr = null;
    var fixture, pan_offset, address = null;

    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            f_selector = 'f_bus-' + bus.id;
            $('#'+f_selector+' :selected').each(function(i, selected) {
                tokens = selected.value.split('-');
                f_addr = parseInt(tokens[tokens.length - 1]);
                fixture = get_fixture_byaddr(bus.id, f_addr);
                pan_offset = fixture.channels.pan;
                address = f_addr + pan_offset;
                bus.buffer[address] = value;
            });

            transfer_buffer(bus);
        }
    }

}


function set_fixture_tilt(value) {
    var bus_name, bus, f_selector, tokens, f_addr = null;
    var fixture, pan_offset, address = null;

    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            f_selector = 'f_bus-' + bus.id;
            $('#'+f_selector+' :selected').each(function(i, selected) {
                tokens = selected.value.split('-');
                f_addr = parseInt(tokens[tokens.length - 1]);
                fixture = get_fixture_byaddr(bus.id, f_addr);
                pan_offset = fixture.channels.tilt;
                address = f_addr + pan_offset;
                bus.buffer[address] = value;
            });

            transfer_buffer(bus);
        }
    }
}


function setup_colorwheel() {
    color_wheel = Raphael.colorwheel($("#color_wheel"), 300, 180);

    var onchange_label = $(".cw_onchange");
    var ondrag_label = $(".cw_ondrag");
    color_wheel.color("#F00");

    color_wheel.ondrag(
        function() {
            ondrag_label.show();
        },
        function() {
            ondrag_label.hide();
        }
    );

    color_wheel.onchange(function(color) {
        set_fixture_color(color);
    });
}


function setup_pantilt_sliders() {
    $('#pan').slider({
        height: 300,
        formatter: function(value) {
            return 'Current value: ' + value;
        }
    });
    $('#pan').on('slide', function(e) {
        e.preventDefault();
        set_fixture_pan(e.value);
    });


    $('#tilt').slider({
        height: 300,
        formatter: function(value) {
            return 'Current value: ' + value;
        }
    });
    $('#tilt').on('slide', function(e) {
        set_fixture_tilt(e.value);
    });
}


/* Various functions used to render pages
 */

function render_direct() {
    /* Function which returns a rendered view of the direct page
     */
    var content, bus_name, bus, f_name, fixture, opt_name = null;

    content = '<div class="row">';

    // Column containing bus/fixture selection widgets
    content += '<div class="col-md-3">';
    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            content += '<h4>' + bus.name + '</h4>';

            content += '<select multiple id="f_bus-' + bus.id + '" class="form-control">';
            for (f_name in dmx_config[bus_name].fixtures) {
                fixture = dmx_config[bus_name].fixtures[f_name];
                opt_name = 'f_bus-' + bus.id+ '-' + fixture.address;
                if (dmx_config[bus_name].fixtures.hasOwnProperty(f_name)) {
                    content += '<option value="'+ opt_name + '">' + f_name + '</option>';
                }
            }
            content += '</select>';
        }
    }
    content += '</div>';

    // Column containing the fixture controls
    content += '<div class="col-md-4">';
    content += '<div id="color_wheel" style="float:left; margin-right:20px; width:300px; text-align:left;"></div>';
    content += '</div>';
    content += '<div class="col-md-1" style="padding-top:50px; padding-left:0px;">';
    content += '<input id="tilt" data-slider-id="tilt_slider" type="text" data-slider-min="0" data-slider-max="255" data-slider-step="4" data-slider-value="128" data-slider-orientation="vertical"/>';
    content += '</div>';

    content += '<div class="row">';
    content += '<div class="col-md-4 col-md-offset-3" style="padding-left:60px; padding-top:20px">';
    content += '<input id="pan" data-slider-id="pan_slider" type="text" data-slider-min="0" data-slider-max="255" data-slider-step="4" data-slider-value="128"/>';
    content += '</div>';

    content += '</div>';

    return content;
}


function render_program() {
    /* Function which returns a rendered view of the programming page
     */
    var content = null;

    content = 'view_program';
    return content;
}

/* Various functions used to view pages
 */

function view_direct() {
    /* Displays and activates the direct control view
     */
    var content = null;
    
    content = render_direct();
    $('#content').html(content);
    setup_colorwheel();
    setup_pantilt_sliders();
}


function view_program() {
    /* Displays and actives the programming view
     */
    var content = null;

    content = render_program();
    $('#content').html(content);
}


function setup_main_eventhandlers() {
    /* Setup eventhandlers for the main window
     */
    $('#v_direct').click(function(e) {
        e.preventDefault();
        view_direct();
    });

    $('#v_program').click(function(e) {
        e.preventDefault();
        view_program();
    });
}


function main() {
    /* Main program to run
     */

    $(document).ready(function() {
        dmx_config = get_json('/v1/config');
        setup_main_eventhandlers();
        view_direct();
    });
}

// Call main function
main();
