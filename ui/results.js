
    var MAX_GROUPS=3
    
    function setEvents() {
        groups.addEventListener("click", function(e){
            // e.target refers to the clicked element
            if (e.target.className.includes("moveable")) {
                let pid = e.target.attributes.pid.nodeValue;
                let gid = parseInt(e.target.attributes.gid.nodeValue);
                let itemMidWidth = e.target.offsetWidth/2;
                if (e.offsetX < itemMidWidth) {
                    movePupil(pid, gid, gid+1);
                } else {
                    if (pid > 0)
                        movePupil(pid, gid, gid-1);
                }
            }
        });
    }

    function movePupil(pid, gidSource, gidTarget) {
        
        if (gidTarget < MAX_GROUPS && gidTarget >=0) {
            var r = confirm("להעביר תלמיד?");
            if (r == true) {
                console.log("מעביר");
            } 
        }
    }

