(function() {
    var state = {
        lat:0, lng:0, alt:null, acc:null, altacc:null, hdg:null, spd:null,
        alpha:0, beta:0, gamma:0,
        bounds:null,
        altbounds:null,
        posts:null,
        interpolate:[],
        pending:false,
        stale:false,
        initialized:false,
        watchId:null,
        scene:null,
        camera:null,
        renderer:null,
        scenePosts:[],
        inputVisible:false,
        aboutVisible:false,
    };

    function initialize() {
        state.initialized = true;
        state.stale = true;
        $("#unsupported").hide();
        window.addEventListener("deviceorientation", updateOrientation, false);
        state.renderer = new THREE.CSS3DRenderer();
        state.camera = new THREE.PerspectiveCamera(45, window.innerWidth/window.innerHeight, 50, 1000);
        state.scene = new THREE.Scene();
        state.scene.add(state.camera);
        state.renderer.setSize(window.innerWidth, window.innerHeight);
        state.renderer.domElement.style.position = "absolute";
        $("#viewport").append(state.renderer.domElement);
        $(state.renderer.domElement).click(function() {
            if (!state.inputVisible) {
                state.inputVisible = true;
                $("#input").fadeIn();
                if (state.aboutVisible) {
                    state.aboutVisible = false;
                    $("#about").fadeOut();
                }
            } else if (!state.aboutVisible) {
                state.inputVisible = false;
                $("#input").fadeOut();
                state.aboutVisible = true;
                $("#about").fadeIn();
            } else {
                state.inputVisible = false;
                $("#input").fadeOut();
                state.aboutVisible = false;
                $("#about").fadeOut();
            }
        });
        window.addEventListener("resize", updateSize, false);
        $.ajax({
            type:"GET",
            url:"/v1/version",
            dataType:"text",
            success:function(data, status) {
                $("#version").text("version " + data);
                setTimeout(function() { $("#about").fadeOut(); }, 2000);
            },
            error:function(xhr, status) {
                setTimeout(function() { $("#about").fadeOut(); }, 2000);
            },
        });
    }

    function showPos(status) {
        if (!state.initialized)
            initialize();
        $("#status").text(status+","+(state.stale && !state.pending));
        $("#lat").text(state.lat);
        $("#lng").text(state.lng);
        $("#alt").text(state.alt);
        $("#acc").text(state.acc);
        $("#altacc").text(state.altacc);
        $("#hdg").text(state.hdg);
        $("#spd").text(state.spd);
        $("#alpha").text(state.alpha);
        $("#beta").text(state.beta);
        $("#gamma").text(state.gamma);
        $("#bounds").text(state.bounds+","+state.altbounds);
        $("#posts").text(JSON.stringify(state.posts));
        if (!state.bounds || state.lat < state.bounds[0] || state.lat > state.bounds[2] || state.lng < state.bounds[1] || state.lng > state.bounds[3]) {
            state.stale = true;
        }
        if (state.stale && !state.pending) {
            state.stale = false;
            state.pending = true;
            $.ajax({
                type:"POST",
                url:"/v1/get",
                data:"lat="+state.lat+"&lng="+state.lng+"&alt="+state.alt+"&acc="+state.acc+"&altacc="+state.altacc+"&hdg="+state.hdg+"&spd="+state.spd+"&alpha="+state.alpha+"&beta="+state.beta+"&gamma="+state.gamma,
                dataType:"json",
                success:function(data, status) {
                    state.bounds = data.bounds;
                    state.posts = data.posts;
                    state.pending = false;
                    state.altbounds = getAltbounds();
                    showPos("post");
                    updateScene();
                },
                error:function(xhdr, status) {
                    state.pending = false;
                },
            });
        }
    }

    function getAltbounds() {
        var altbounds = [state.alt-1,state.alt+1];
        for (var i = 0; i < state.posts.length; i++) {
            if (state.posts[i].loc.alt < altbounds[0])
                altbounds[0] = state.posts[i].loc.alt;
            else if (state.posts[i].loc.alt > altbounds[1])
                altbounds[1] = state.posts[i].loc.alt;
        }
        return altbounds;
    }

    function updatePos(pos) {
        var dlat = pos.coords.latitude - state.lat;
        var dlng = pos.coords.longitude - state.lng;
        var dalt = pos.coords.altitude - state.alt;
        state.interpolate = [];
        for (var i = 0; i < 20; i++) {
            var p = {
                lat: state.lat + i*dlat/20,
                lng: state.lng + i*dlng/20,
                alt: state.alt + i*dalt/20,
            };
            state.interpolate.push(p);
        }
        state.lat = pos.coords.latitude;
        state.lng = pos.coords.longitude;
        state.alt = pos.coords.altitude;
        state.acc = pos.coords.accuracy;
        state.altacc = pos.coords.altitudeAccuracy;
        state.hdg = pos.coords.heading;
        state.spd = pos.coords.speed;
        showPos("pos="+pos);
        animate();
        $("#camerapos").text(state.camera.position.x+","+state.camera.position.y+","+state.camera.position.z);
    }

    function animate() {
        if (state.interpolate.length == 0) {
            setPosition(state.camera, state);
            setDirection(state.camera, state, 1);
            state.renderer.render(state.scene, state.camera);
            return;
        }
        var pos = state.interpolate.shift();
        setPosition(state.camera, pos);
        setDirection(state.camera, state, 1);
        state.renderer.render(state.scene, state.camera);
        requestAnimationFrame(animate);
    }

    function setDirection(obj, loc, sign) {
        obj.lookAt(new THREE.Vector3(obj.position.x + sign*Math.sin(loc.alpha*Math.PI/180), obj.position.y + sign*0, obj.position.z + sign*Math.cos(loc.alpha*Math.PI/180)));
    }

    function setPosition(obj, loc) {
        obj.position.x = ((loc.lat - state.bounds[0])/(state.bounds[2] - state.bounds[0]) - 0.5)*1000;
        if (typeof(loc.alt) == "number")
            obj.position.y = ((loc.alt - state.altbounds[0])/(state.altbounds[1] - state.altbounds[0]) - 0.5)*100;
        else
            obj.position.y = 0;
        obj.position.z = ((loc.lng - state.bounds[1])/(state.bounds[3] - state.bounds[1]) - 0.5)*1000;
    }

    function dangle(a1, a2) {
        var d = a1 - a2;
        while (d > 180)
            d -= 360;
        while (d <= -180)
            d += 360;
        return d;
    }

    function updateOrientation(orient) {
        if (state.stale || Math.abs(dangle(orient.alpha, state.alpha)) > 1 || Math.abs(dangle(orient.beta, state.beta)) > 1 || Math.abs(dangle(orient.gamma, state.gamma)) > 1) {
            state.alpha = orient.alpha;
            state.beta = orient.beta;
            state.gamma = orient.gamma;
            showPos("orient="+orient);
            setDirection(state.camera, state, 1);
            state.renderer.render(state.scene, state.camera);
        }
    }

   function updateScene() {
        for (var i = 0; i < state.scenePosts.length; i++)
            state.scene.remove(state.scenePosts[i]);
        state.scenePosts = [];
        for (var i = 0; i < state.posts.length; i++) {
            var div = document.createElement("div");
            var lines = state.posts[i].msg.split("\n");
            for (var j = 0; j < lines.length; j++) {
                if (j > 0) {
                    $(div).append(document.createElement("br"));
                }
                $(div).append(document.createTextNode(lines[j]));
            }
            var scenePost = new THREE.CSS3DObject(div);
            setPosition(scenePost, state.posts[i].loc);
            setDirection(scenePost, state.posts[i].loc, -1);
            state.scenePosts.push(scenePost);
            state.scene.add(scenePost);
        }
        state.renderer.render(state.scene, state.camera);
    }

    function updateSize() {
        state.camera.aspect = window.innerWidth/window.innerHeight;
        state.camera.updateProjectionMatrix();
        state.renderer.setSize(window.innerWidth, window.innerHeight);
        state.renderer.render(state.scene, state.camera);
    }

    $(document).ready(function() {
        $("#post").click(function() {
            var msg = $("#msg").val();
            $("#msg").val(null);
            state.inputVisible = false;
            $("#input").fadeOut();
            $.ajax({
                type:"POST",
                url:"/v1/post",
                data:"lat="+state.lat+"&lng="+state.lng+"&alt="+state.alt+"&acc="+state.acc+"&altacc="+state.altacc+"&hdg="+state.hdg+"&spd="+state.spd+"&alpha="+state.alpha+"&beta="+state.beta+"&gamma="+state.gamma+"&msg="+encodeURIComponent(msg),
                success:function(data, status) {
                    state.stale = true;
                    showPos("post");
                },
            });
        });

        $("#cancel").click(function() {
            state.inputVisible = false;
            $("#input").fadeOut();
        });

        state.watchId = navigator.geolocation && navigator.geolocation.watchPosition(updatePos);
    });
})()
