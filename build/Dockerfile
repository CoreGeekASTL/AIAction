FROM {{csp-go-image}}

USER root
RUN mkdir -p /opt/csp/gids/module
COPY . /opt/csp/gids/module
RUN cat /opt/csp/gids/module/sudoer.config >> /etc/sudoers && \
    rm -rf /opt/csp/gids/module/sudoer.config && \
    rm -rf /opt/csp/gids/module/Dockerfile && \
    chown paas:paas -R /opt/csp/gids && \
    chmod +x /opt/csp/gids/module/gids && \
    chmod +x /opt/csp/gids/module/*.sh
USER paas