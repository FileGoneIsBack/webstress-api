$(document).ready(function() {
    // Initialize Particles.js
    if (typeof particlesJS !== 'undefined') {
      particlesJS('particles-js', {
        particles: {
          number: {
            value: 80,
            density: {
              enable: true,
              value_area: 800
            }
          },
          color: {
            value: "#ea384c"
          },
          shape: {
            type: "circle",
            stroke: {
              width: 0,
              color: "#000000"
            },
            polygon: {
              nb_sides: 5
            }
          },
          opacity: {
            value: 0.5,
            random: false,
            anim: {
              enable: false,
              speed: 1,
              opacity_min: 0.1,
              sync: false
            }
          },
          size: {
            value: 3,
            random: true,
            anim: {
              enable: false,
              speed: 40,
              size_min: 0.1,
              sync: false
            }
          },
          line_linked: {
            enable: true,
            distance: 150,
            color: "#ea384c",
            opacity: 0.4,
            width: 1
          },
          move: {
            enable: true,
            speed: 6,
            direction: "none",
            random: false,
            straight: false,
            out_mode: "out",
            bounce: false,
            attract: {
              enable: false,
              rotateX: 600,
              rotateY: 1200
            }
          }
        },
        interactivity: {
          detect_on: "canvas",
          events: {
            onhover: {
              enable: true,
              mode: "repulse"
            },
            onclick: {
              enable: true,
              mode: "push"
            },
            resize: true
          },
          modes: {
            grab: {
              distance: 400,
              line_linked: {
                opacity: 1
              }
            },
            bubble: {
              distance: 400,
              size: 40,
              duration: 2,
              opacity: 8,
              speed: 3
            },
            repulse: {
              distance: 200,
              duration: 0.4
            },
            push: {
              particles_nb: 4
            },
            remove: {
              particles_nb: 2
            }
          }
        },
        retina_detect: true
      });
    }
  
    // Navbar scroll effect
    $(window).scroll(function() {
      if ($(window).scrollTop() > 50) {
        $('#main-nav').addClass('scrolled');
      } else {
        $('#main-nav').removeClass('scrolled');
      }
    });
  
    // Mobile menu toggle
    $('.mobile-menu-btn').click(function() {
      $('.nav-links').slideToggle();
    });
  
    // Pricing tabs
    $('.tab-btn').click(function() {
      const tabId = $(this).data('tab');
      
      // Update buttons
      $('.tab-btn').removeClass('active');
      $(this).addClass('active');
      
      // Update content
      $('.tab-pane').removeClass('active');
      $('#' + tabId).addClass('active');
    });
  
    // Smooth scroll for navigation links
    $('a[href^="#"]').on('click', function(e) {
      e.preventDefault();
      
      const target = this.hash;
      const $target = $(target);
      
      $('html, body').animate({
        'scrollTop': $target.offset().top - 70
      }, 1000, 'swing');
    });
  });
  