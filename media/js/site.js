// BEGIN: LocalStorage functions
function store_get(k) {
    try {
        v = s.get(k)
        console.log('store_get: '+k+' = '+v)
        if (v == 'false') {
            return false
        } else if (v == 'true') {
            return true
        }
        return JSON.parse(v)
    } catch (err) {
        console.log('store_get: could not retrieve '+k)
        return null
    }
}

function store_set(k, v) {
    console.log('store_set: '+k+' = '+JSON.stringify(v))
    s.set(k, JSON.stringify(v))
    s.save()
}

function reset_store() {
    console.log('resetting LocalStorage')
    store_set('spots', {})
    store_set('groups', {})
}


// BEGIN: Application initialization
function load_configuration() {
    $.ajax({
        url: '/config',
        type: 'get',
        success: function(response) {
            store_set('spots', response['spots'])
            load_spots(response['spots'])

            store_set('groups', response['groups'])
            load_groups(response['groups'])
        },
        error: function(xhr, textStatus, errorThrown) {
            console.log('Failed to load configuration: '+textStatus+', '+errorThrown)
        }
    })
}

function load_spots(spots) {
    var content = '<h4>Available spots</h4>'
    content += '<div class="row">'
    content += '<div class="span1"><b>Name</b></div>'
    content += '<div class="span3"><b>Description</b></div>'
    content += '<div class="span2"><b>Location</b></div>'
    content += '<div class="span1"><b>Slot</b></div>'
    content += '<div class="span1"><b>Color</b></div>'
    content += '</div>'

    var keys = Object.keys(spots)
    keys.sort()
    for (var i=0; i<keys.length; i++) {
        var spot_name = keys[i]
        var spot_info = spots[spot_name]
        content += '<div class="row">'
        content += '<div class="span1">'+spot_name+'</div>'
        content += '<div class="span3">'+spot_info['description']+'</div>'
        content += '<div class="span2">'+spot_info['location']+'</div>'
        content += '<div class="span1">'+spot_info['slot']+'</div>'
        content += '<div class="span1">#000000</div>'
        content += '</div>'
    }

    $('#available_spots').html(content)
}

function load_groups(groups) {
    var content = '<h4>Available groups</h4>'
    content += '<div class="row">'
    content += '<div class="span1"><b>Name</b></div>'
    content += '<div class="span3"><b>Description</b></div>'
    content += '<div class="span2"><b>Location</b></div>'
    content += '<div class="span3"><b>Spots</b></div>'
    content += '</div>'

    var keys = Object.keys(groups)
    keys.sort()
    for (var i=0; i<keys.length; i++) {
        var group_name = keys[i]
        var group_info = groups[group_name]
        content += '<div class="row">'
        content += '<div class="span1">'+group_name+'</div>'
        content += '<div class="span3">'+group_info['description']+'</div>'
        content += '<div class="span2">'+group_info['location']+'</div>'
        content += '<div class="span3">'+group_info['spots']+'</div>'
        content += '</div>'
    }

    $('#available_groups').html(content)
}

function setup_colorwheel() {
    var cw = Raphael.colorwheel($(".cw_id")[0],300)
    var onchange_el = $(".cw_onchange")
    var ondrag_el = $(".cw_ondrag")
    cw.color("#F00");

    function start() {
        ondrag_el.show()
    }

    function stop() {
        ondrag_el.hide()
    }

    cw.ondrag(start, stop);

    cw.onchange(function(color) {
        var colors = [
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
            parseInt(color.r), parseInt(color.g), parseInt(color.b),
        ]
        var frame_data = {}
        frame_data['frame'] = colors
        frame_data['duration'] = 20

        $.ajax({
            url: '/frame',
            type: 'put',
            data: JSON.stringify(frame_data),
            dataType: 'json',
            error: function(xhr, textStatus, errorThrown) {
                console.log('request failed: '+textStatus+'; '+errorThrown)
            }
        })
    })
}

function setup_application() {
    $('#overview a').click(function (e) {
        e.preventdefault();
        $(this).tab('show');
    })

    $('#colorwheel a').click(function (e) {
        e.preventdefault();
        $(this).tab('show');
    })
}

function main() {
    $(document).ready(function() {
        s = new Persist.Store('inti')
        reset_store()
        load_configuration()
        setup_application()
        setup_colorwheel()
    })
}

/* GO GO GO */
main()
