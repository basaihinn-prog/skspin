<template>
   <div class="about-container">
    <div class="about-header">
      <h1>About Us</h1>
    </div>

    <div class="about-content">
      <h2>At SkyXSpin, we specialize in innovative slot replication</h2>
      <p>
        Pushing the boundaries of technology to deliver an unparalleled gaming experience. We pride ourselves on leading the way in server efficiency and technical excellence, ensuring that players enjoy high-quality gaming content at exceptional speeds.
      </p>
      <p>
        Our process involves taking the original front-end and source code files and compiling them into a lightweight framework. This approach ensures zero-latency and maintains the premium quality expected from the leading slot providers in the Asian market. With a focus on efficiency, precision, and reliability, our platform guarantees players an unmatched gaming experience and next to no loading times.
      </p>
    </div>

    <div class="about-section">
      <h3>Our Remote Game Server</h3>
      <ul>
        <li><strong>Seamless Integration:</strong> Experience effortless connectivity via our dedicated delivery network, providing operators and players access to the industry’s top-rated slot games.</li>
        <li><strong>Lowest GGR Rates:</strong> Maximize your revenue potential with the most competitive Gross Gaming Revenue rates in the market.</li>
        <li><strong>Complete RTP Control:</strong> Adjust Return-to-Player percentages (RTP) to match your desired player experience, giving you ultimate flexibility in game management from your own Admin panel.</li>
        <li><strong>Reliability and Speed:</strong> Our platform ensures consistent, fast delivery of slot content, maintaining uptime and performance across all devices.</li>
      </ul>
    </div>

    <div class="about-section">
      <h3>Why Choose SkyXSpin?</h3>
      <ul>
        <li><strong>Comprehensive Game Portfolio:</strong> Access the most popular slots, including mini games, tailored to the preferences of your market.</li>
        <li><strong>Innovative Game Replication:</strong> Utilizing advanced technologies, we replicate slot games to ensure flawless execution and top-tier quality.</li>
        <li><strong>Customizable Solutions:</strong> From game configurations to branding, we provide tailored options to fit your business needs.</li>
        <li><strong>Player-Centric Approach:</strong> Our games are designed to captivate players, keeping engagement high and driving long-term loyalty.</li>
      </ul>
    </div>

    <div class="about-footer">
      <p>
        At SkyXSpin, we are more than just a slot clone provider. We are innovators and creators, dedicated to setting new standards in online gaming.
      </p>
      <p>
        Discover the future of gaming with SkyXSpin—where speed, quality, and innovation converge to redefine the slot experience.
      </p>
    </div>
   </div>
    
</template>


<style scoped>
.bottom_text{
    margin-top: 5rem;
}

.about-container {
  max-width: 93.75rem; /* 1500px */
  margin: 0 auto;
  padding: 2.5rem; /* 40px */
}
@media only screen and (max-width:768px){
    .about-container{
        padding-top: 5rem;
        text-align: center; 
    }

}

.about-header {
  margin-bottom: 1.875rem; /* 30px */
  text-align: center;
}

.about-header h1 {
  font-size: 4.125rem; /* 50px */
  color: #ffffff;
  font-weight: bold;
}

.about-content h2{
  margin-bottom: 2.5rem; /* 40px */
  color: #ffffff;
}

.about-content h2 {
  font-size: 1.75rem; /* 28px */
  margin-bottom: 1.25rem; /* 20px */
}


.about-content p {
  font-size: 1.125rem; /* 18px */
  line-height: 1.6;
  margin-bottom: 1.25rem; /* 20px */
}


.about-section {
  margin-bottom: 2.5rem; /* 40px */
}

.about-section h3 {
  font-size: 1.5rem; /* 24px */
  margin-bottom: 1.25rem; /* 20px */
  color: #ffffff;
}

.about-section ul {
  list-style-type: none;
  padding-left: 0;
  font-size: 1.125rem; /* 18px */
}

.about-section li {
  margin-bottom: 0.625rem; /* 10px */
}

.about-footer {
  margin-top: 2.5rem; /* 40px */
  font-size: 1.125rem; /* 18px */
  color: #ffffff;
}

@media only screen and (max-width:768px){
    .about-content p, .about-section ul {
        font-size: 1rem; /* 18px */
        line-height: 1.6;
    }   
}

@media only screen and (max-width:768px){
    .about-content h2 {
        font-size: 1.55rem;
    }   
    .about-header h1{
        font-size: 3.125rem
    }
}
</style>



<script lang="ts">

const hrefText = window.location.href

export default {
  name: 'AboutUs',
};
</script>