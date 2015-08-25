/*jslint browser: true*/
/*global $, jQuery, alert*/
"use strict";

var dmx_config = new Array();


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


/* Various functions used to render pages
 */

function render_direct() {
    /* Function which returns a rendered view of the direct page
     */
    var content, bus_name, bus, f_name, fixture = null;

    content = '<div class="row">';
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
