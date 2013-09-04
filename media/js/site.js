var groups = []
var controllers = []
var current_target

// BEGIN: LocalStorage functions
function store_get(k) {
    try {
        v = s.get(k)
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
    s.set(k, JSON.stringify(v))
    s.save()
}

function reset_store() {
    console.log('resetting LocalStorage')
    store_set('spots', {})
    store_set('groups', {})
    store_set('current_target', null)
}

function rgb_to_hex(r, g, b) {
    return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1,7);
}

// BEGIN: Application initialization
function load_configuration() {
    $.ajax({
        url: '/config',
        type: 'get',
        success: function(response) {
            controllers = response['Controllers']
            load_controllers(response['Controllers'])

            groups = response['Groups']
            load_groups(response['Groups'])

            setup_target_menu(response['Controllers'], response['Groups'])

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
            var color = rgb_to_hex(spot_info['R'],spot_info['G'],spot_info['B'])
            content += '<div class="row">'
            content += '<div class="span2">'+spot_info['Name']+'</div>'
            content += '<div class="span3">'+spot_info['Description']+'</div>'
            content += '<div class="span1" style="background-color:'+color+'">'+color+'</div>'
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
            var color = rgb_to_hex(spot_info['R'],spot_info['G'],spot_info['B'])
            content += '<div class="row">'
            content += '<div class="span2">'+spot_info['Name']+'</div>'
            content += '<div class="span3">'+spot_info['Description']+'</div>'
            content += '<div class="span1" style="background-color:'+color+'">'+color+'</div>'
            content += '</div>'
        }
        content += '</div>'
        content += '</div>'
    }

    $('#available_groups').html(content)
}

function setup_target_menu(controllers, groups) {
    var content = '<ul class="nav nav-list well">'
    content += '<li class="nav-header">Select target</li>'
    content += '<ul class="nav nav-list">'
    content += '<li class="nav-header">By group</li>'
    for (gid = 0; gid < groups.length; gid++) {
        content += '<li id="m_group_'+gid+'"><a href="#">'+groups[gid].Name+'</a></li>'
    }
    content += '</ul>'
    content += '<ul class="nav nav-list">'
    content += '<li class="nav-header">By controller</li>'

    for (cid = 0; cid < controllers.length; cid++) {
        content += '<li id="m_controller_'+cid+'"><a href="#">'+controllers[cid].Name+'</a></li>'
    }
    content += '</ul>'
    content += '</ul>'

    $('#cw_menu').html(content)

    for (gid = 0; gid < groups.length; gid++) {
        $('#m_group_'+gid).click(function(e) {
            $('#'+current_target).removeClass('active')
            $(this).addClass('active')
            current_target = this.id
        })
    }

    for (cid = 0; cid < controllers.length; cid++) {
        $('#m_controller_'+cid).click(function(e) {
            $('#'+current_target).removeClass('active')
            $(this).addClass('active')
            current_target = this.id
        })
    }
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
        console.log(current_target)

        if (current_target.indexOf('m_group_') > -1) {
            id = parseInt(current_target.replace('m_group_', ''))
            target = groups[id]
        } else {
            id = parseInt(current_target.replace('m_controller_', ''))
            target = controllers[id]
        }

        var colors = []
        for (i = 0; i < target.BufSize+1; i++) {
            colors[i] = 0
        }

        for (sid = 0; sid < target.Spots.length; sid++) {
            offset = target.Spots[sid].Id - 1
            console.log(offset)
            colors[offset] = parseInt(color.r)
            colors[offset+1] = parseInt(color.g)
            colors[offset+2] = parseInt(color.b)
        }

        var frame_data = {}
        frame_data['frame'] = colors
        frame_data['duration'] = 20
        console.log(colors)

        $.ajax({
            url: '/frame',
            type: 'put',
            data: JSON.stringify(frame_data),
            ataType: 'json',
            error: function(xhr, textStatus, errorThrown) {
                console.log('request failed: '+textStatus+'; '+errorThrown)
            }
        })
    })
}

function setup_eventhandlers() {
    $('#a_overview').click(function(e) {
        console.log("a_overview clicked")
        load_configuration()
    })
}

function main() {
    $(document).ready(function() {
        s = new Persist.Store('inti')
        reset_store()
        setup_eventhandlers()
        load_configuration()
        setup_colorwheel()
    })
}

/* GO GO GO */
main()
