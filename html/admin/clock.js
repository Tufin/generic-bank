function initClock(id, time) {


    var clock = new function() {
        canvas = document.getElementById(id);

        this.ctx = canvas.getContext("2d");
        this.radius = canvas.height / 2;
    }

    clock.drawFace = function() {
        var grad;

        ctx = this.ctx;
        radius = this.radius;

        ctx.beginPath();
        ctx.arc(0, 0, radius, 0, 2*Math.PI);
        ctx.fillStyle = 'white';
        ctx.fill();

        grad = ctx.createRadialGradient(0,0,radius*0.95, 0,0,radius*1.05);
        grad.addColorStop(0, '#333');
        grad.addColorStop(0.5, 'white');
        grad.addColorStop(1, '#333');
        ctx.strokeStyle = grad;
        ctx.lineWidth = radius*0.1;
        ctx.stroke();

        ctx.beginPath();
        ctx.arc(0, 0, radius*0.1, 0, 2*Math.PI);
        ctx.fillStyle = '#333';
        ctx.fill();
    }

    clock.drawNumbers = function() {
        var ang;
        var num;
        ctx.font = radius*0.15 + "px arial";
        ctx.textBaseline="middle";
        ctx.textAlign="center";
        for(num= 1; num < 13; num++){
            ang = num * Math.PI / 6;
            ctx.rotate(ang);
            ctx.translate(0, -radius*0.85);
            ctx.rotate(-ang);
            ctx.fillText(num.toString(), 0, 0);
            ctx.rotate(ang);
            ctx.translate(0, radius*0.85);
            ctx.rotate(-ang);
        }
    }

    clock.drawHand = function(pos, length, width) {
        ctx.beginPath();
        ctx.lineWidth = width;
        ctx.lineCap = "round";
        ctx.moveTo(0,0);
        ctx.rotate(pos);
        ctx.lineTo(0, -length);
        ctx.stroke();
        ctx.rotate(-pos);
    }

    clock.drawTime = function() {
        var now = new Date(time);
        var hour = now.getHours();
        var minute = now.getMinutes();
        var second = now.getSeconds();
        //hour
        hour=hour%12;
        hour=(hour*Math.PI/6)+(minute*Math.PI/(6*60))+(second*Math.PI/(360*60));
        clock.drawHand(hour, radius*0.5, radius*0.07);
        //minute
        minute=(minute*Math.PI/30)+(second*Math.PI/(30*60));
        clock.drawHand(minute, radius*0.8, radius*0.07);
        // second
        second=(second*Math.PI/30);
        clock.drawHand(second, radius*0.9, radius*0.02);
    }

    clock.drawClock = function() {
        clock.drawFace();
        clock.drawNumbers();
        clock.drawTime();
        clock.drawHand();
    };

    clock.ctx.translate(clock.radius, clock.radius);
    clock.radius = clock.radius * 0.90;
    setInterval(clock.drawClock, 1000);

}


