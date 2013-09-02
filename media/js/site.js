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
            store_set('controllers', response['Controllers'])
            load_controllers(response['Controllers'])

            store_set('groups', response['Groups'])
            load_groups(response['Groups'])
        },
        error: function(xhr, textStatus, errorThrown) {
            console.log('Failed to load configuration: '+textStatus+', '+errorThrown)
        }
    })
}

function load_controllers(controllers) {
    var content = '<h4>Available Controllers</h4>'

    for (cid=0; cid < controllers.length; cid++) {
        content += '<div class="row">'
        content += '<div class="span4">'
        content += 'Name: '+controllers[cid].Name+'<br/>'
        content += 'Description: '+controllers[cid].Description+'<br/>'
        content += '</div>'

        content += '<div class="span8">'
        for (var sid=0; sid < controllers[cid].Spots.length; sid++) {
            var spot_info = controllers[cid].Spots[sid]
            content += '<div class="row">'
            content += '<div class="span2">'+spot_info['Name']+'</div>'
            content += '<div class="span3">'+spot_info['Description']+'</div>'
            content += '<div class="span1">'+spot_info['R']+'</div>'
            content += '<div class="span1">'+spot_info['G']+'</div>'
            content += '<div class="span1">'+spot_info['B']+'</div>'
            content += '</div>'
        }
        content += '</div>'
        content += '</div>'
    }

    $('#available_spots').html(content)
}

function load_groups(groups) {
    var content = '<h4>Available groups</h4>'

    for (gid=0; gid < groups.length; gid++) {
        content += '<div class="row">'
        content += '<div class="span4">'
        content += 'Name: '+groups[gid].Name+'<br/>'
        content += 'Description: '+groups[gid].Description+'<br/>'
        content += '</div>'

        content += '<div class="span8">'
        for (var sid=0; sid < groups[gid].Spots.length; sid++) {
            var spot_info = groups[gid].Spots[sid]
            content += '<div class="row">'
            content += '<div class="span2">'+spot_info['Name']+'</div>'
            content += '<div class="span3">'+spot_info['Description']+'</div>'
            content += '<div class="span1">'+spot_info['R']+'</div>'
            content += '<div class="span1">'+spot_info['G']+'</div>'
            content += '<div class="span1">'+spot_info['B']+'</div>'
            content += '</div>'
        }
        content += '</div>'
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
