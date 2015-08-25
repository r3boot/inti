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


function set_fixture_color(color) {
    console.log("Setting fixture color");
    console.log(color);
}


function set_fixture_pan(value) {
    console.log("Setting fixture pan");
    console.log(value);
}


function set_fixture_tilt(value) {
    console.log("Setting fixture tilt");
    console.log(value);
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
    var content, bus_name, bus, f_name, fixture = null;

    content = '<div class="row">';

    // Column containing bus/fixture selection widgets
    content += '<div class="col-md-3">';
    for (bus_name in dmx_config) {
        if (dmx_config.hasOwnProperty(bus_name)) {
            bus = dmx_config[bus_name];
            content += '<h4>' + bus.name + '</h4>';

            content += '<select multiple class="form-control">';
            for (f_name in dmx_config[bus_name].fixtures) {
                if (dmx_config[bus_name].fixtures.hasOwnProperty(f_name)) {
                    content += '<option value="'+ f_name + '">' + f_name + '</option>';
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
    content += '<input id="tilt" data-slider-id="tilt_slider" type="text" data-slider-min="0" data-slider-max="255" data-slider-step="1" data-slider-value="128" data-slider-orientation="vertical"/>';
    content += '</div>';

    content += '<div class="row">';
    content += '<div class="col-md-4 col-md-offset-3" style="padding-left:60px; padding-top:20px">';
    content += '<input id="pan" data-slider-id="pan_slider" type="text" data-slider-min="0" data-slider-max="255" data-slider-step="1" data-slider-value="128"/>';
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
