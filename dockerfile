# Due to an unkown bug on loading dependencies into the container, image uses the binary rather than go files

# Use a base image
FROM debian:latest

# Create directory
RUN mkdir -p /usr/local/bin

# Copy the 'ckb' binary into the container
COPY ckb /usr/local/bin/

# Set the working directory
WORKDIR /usr/local/bin/

COPY ./templates/ ./templates/

# Make 'ckb' executable
RUN chmod +x ckb

# Expose port to access app
EXPOSE 8081

# Define the command to run when the container starts
CMD ["ckb"]