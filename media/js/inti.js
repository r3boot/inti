/*jslint browser: true*/
/*global $, jQuery, alert*/
"use strict";

/* Various functions used to render pages
 */

function render_direct() {
    /* Function which returns a rendered view of the direct page
     */
    var content = null;

    content = 'view_direct';
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
        $.ajax({
            url: '/v1/config',
            type: 'get',
            success: function(r) {
                console.log(r);
                setup_main_eventhandlers();
                view_direct();
            },
            error: function(xhr, textStatus, errorThrown) {
                console.log('Failed to load configuration: ' + xhr + ', ' + textStatus + ', ' + errorThrown);
                $('#content').html("Failed to load configuration");
            },
        });
    });
}

// Call main function
main();
